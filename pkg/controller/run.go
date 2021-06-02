package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func (ctrl *Controller) Run(ctx context.Context, param Param) error { //nolint:cyclop
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}
	param.Items = cfg.Items
	state := State{}
	if param.StatePath != "" {
		if err := ctrl.readState(param.StatePath, &state); err != nil {
			return fmt.Errorf("read state (state path: %s): %w", param.StatePath, err)
		}
	} else {
		if err := ctrl.readStateFromCmd(ctx, &state); err != nil {
			return err
		}
	}

	tfPath, err := ctrl.writeTF()
	if tfPath != "" {
		defer os.Remove(tfPath)
	}
	if err != nil {
		return err
	}

	dryRunResult := DryRunResult{}

	for _, rsc := range state.Values.RootModule.Resources {
		if err := ctrl.handleResource(ctx, param, rsc, tfPath, &dryRunResult); err != nil {
			return err
		}
	}
	if param.DryRun {
		if err := yaml.NewEncoder(ctrl.Stdout).Encode(dryRunResult); err != nil {
			return fmt.Errorf("encode dry run result as YAML: %w", err)
		}
	}
	return nil
}

func (ctrl *Controller) readState(statePath string, state *State) error {
	stateFile, err := os.Open(statePath)
	if err != nil {
		return fmt.Errorf("open a state file %s: %w", statePath, err)
	}
	defer stateFile.Close()
	if err := json.NewDecoder(stateFile).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON %s: %w", statePath, err)
	}
	return nil
}

func (ctrl *Controller) readStateFromCmd(ctx context.Context, state *State) error {
	buf := bytes.Buffer{}
	if err := ctrl.tfShow(ctx, &buf); err != nil {
		return err
	}
	if err := json.NewDecoder(&buf).Decode(state); err != nil {
		return fmt.Errorf("parse a state as JSON: %w", err)
	}
	return nil
}

func (ctrl *Controller) writeTF() (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("create a temporal file to write Terraform configuration (.tf): %w", err)
	}
	defer f.Close()
	// read tf from stdin and write a temporal file
	if _, err := io.Copy(f, ctrl.Stdin); err != nil {
		return f.Name(), fmt.Errorf("write standard input to a temporal file %s: %w", f.Name(), err)
	}
	return f.Name(), nil
}

type ResourcePath struct {
	Type string
	Name string
}

func (rp *ResourcePath) Path() string {
	return rp.Type + "." + rp.Name
}

func (ctrl *Controller) handleResource(ctx context.Context, param Param, rsc Resource, hclFilePath string, dryRunResult *DryRunResult) error {
	matched := false
	for _, item := range param.Items {
		f, err := ctrl.handleItem(ctx, rsc, item, hclFilePath, param, dryRunResult)
		if err != nil {
			return fmt.Errorf("handle item (rule: %s): %w", item.Rule.Raw(), err)
		}
		if f {
			matched = true
			break
		}
	}
	if !matched && param.DryRun {
		resourcePath, err := getResourcePath(rsc)
		if err != nil {
			return err
		}
		dryRunResult.NoMatchResources = append(dryRunResult.NoMatchResources, resourcePath.Path())
	}
	return nil
}

func (ctrl *Controller) handleItem(ctx context.Context, rsc Resource, item Item, hclFilePath string, param Param, dryRunResult *DryRunResult) (bool, error) { //nolint:cyclop,funlen
	resourcePath, err := getResourcePath(rsc)
	if err != nil {
		return true, err
	}
	// filter resource by condition
	matched, err := item.Rule.Run(rsc)
	if err != nil {
		return false, fmt.Errorf("check if the rule matches with the resource: %w", err)
	}
	if !matched {
		return false, nil
	}

	if item.Exclude {
		if param.DryRun {
			dryRunResult.ExcludedResources = append(dryRunResult.ExcludedResources, resourcePath.Path())
		}
		return true, nil
	}

	newResourcePath := resourcePath
	if item.ResourceName != nil {
		// compute new resource path
		newResourcePath.Name, err = item.ResourceName.Parse(rsc)
		if err != nil {
			return true, fmt.Errorf("compute a new resource name (template: %s): %w", item.ResourceName.Raw(), err)
		}
	}

	if param.DryRun {
		dryRunResult.MigratedResources = append(dryRunResult.MigratedResources, MigratedResource{
			SourceResourcePath: resourcePath.Path(),
			DestResourcePath:   newResourcePath.Path(),
			TFBasename:         item.TFBasename,
			StateDirname:       item.StateDirname,
			StateBasename:      item.StateBasename,
		})
		return true, nil
	}

	hclFile, err := os.Open(hclFilePath)
	if err != nil {
		return true, fmt.Errorf("open a Terraform configuration %s: %w", hclFilePath, err)
	}
	defer hclFile.Close()

	stateBasename, err := item.StateBasename.Execute(rsc)
	if err != nil {
		return true, fmt.Errorf("render a state_basename: %w", err)
	}

	stateDirname, err := item.StateDirname.Execute(rsc)
	if err != nil {
		return true, fmt.Errorf("render a state_dirname: %w", err)
	}

	tfBasename, err := item.TFBasename.Execute(rsc)
	if err != nil {
		return true, fmt.Errorf("render a tf_basename: %w", err)
	}

	tfPath := filepath.Join(stateDirname, tfBasename)
	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return true, fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	buf := bytes.Buffer{}
	if err := ctrl.getHCL(ctx, resourcePath.Path(), newResourcePath.Path(), hclFile, &buf); err != nil {
		return true, err
	}

	if err := ctrl.stateMv(ctx, filepath.Join(stateDirname, stateBasename), resourcePath.Path(), newResourcePath.Path(), param.SkipState); err != nil {
		return true, err
	}
	// write hcl
	if _, err := io.Copy(tfFile, &buf); err != nil {
		return true, fmt.Errorf("write Terraform configuration to a file %s: %w", tfPath, err)
	}
	return true, nil
}

func (ctrl *Controller) getHCL(
	ctx context.Context, resourcePath, newResourcePath string, hclFile io.Reader, buf io.Writer) error {
	if resourcePath == newResourcePath {
		return ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, buf)
	}
	pp := bytes.Buffer{}
	if err := ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, &pp); err != nil {
		return fmt.Errorf("get a resource from HCL file: %w", err)
	}

	if err := ctrl.blockMv(ctx, "resource."+resourcePath, "resource."+newResourcePath, &pp, buf); err != nil {
		return fmt.Errorf("rename resource: %w", err)
	}
	return nil
}

func getResourcePath(rsc Resource) (ResourcePath, error) { //nolint:unparam
	return ResourcePath{
		Type: rsc.Type,
		Name: rsc.Name,
	}, nil
}
