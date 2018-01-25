package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	// You should never import this, this will try to reinitialize
	// flags, and panic
	//
	// utilflag "k8s.io/apiserver/pkg/util/flag"

	"k8s.io/kubernetes/pkg/kubectl/cmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/resource"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	kresource "k8s.io/kubernetes/pkg/kubectl/resource"
)
import "C"

func translateGetOptions(raw map[string]interface{}) *resource.GetOptions {
	foptions := kresource.FilenameOptions{}

	if fnames, ok := raw["filenames"]; ok {
		foptions.Filenames = fnames.([]string)
	}
	if recursive, ok := raw["recursive"]; ok {
		foptions.Recursive = recursive.(bool)
	}
	result := &resource.GetOptions{
		FilenameOptions: foptions,
	}
	for k, v := range raw {
		switch k {
		case "raw":
			result.Raw = v.(string)
		case "watch":
			result.Watch = v.(bool)
		case "watch_olnly":
			result.WatchOnly = v.(bool)
		case "chunk_size":
			result.ChunkSize = v.(int64)
		case "label_selector":
			result.LabelSelector = v.(string)
		case "field_selector":
			result.FieldSelector = v.(string)
		case "all_namespaces":
			result.AllNamespaces = v.(bool)
		case "namespace":
			result.Namespace = v.(string)
		case "explicit_namespace":
			result.ExplicitNamespace = v.(bool)
		case "ignore_not_found":
			result.IgnoreNotFound = v.(bool)
		case "show_kind":
			result.ShowKind = v.(bool)
		case "export":
			result.Export = v.(bool)
		case "include_uninitialized":
			result.IncludeUninitialized = v.(bool)
		}
	}
	return result
}

func translateCreateOptions(raw map[string]interface{}) *cmd.CreateOptions {
	return &cmd.CreateOptions{}
}

//export ResourceGet
func ResourceGet(optsEncoded string, typeOrName []string) (res string, serr string) {
	defer func() {
		if r := recover(); r != nil {
			serr = fmt.Sprintf("%v\n%s", r, debug.Stack())
		}
	}()
	opts := map[string]interface{}{}
	if err := json.Unmarshal([]byte(optsEncoded), &opts); err != nil {
		return "", err.Error()
	}
	factory := cmdutil.NewFactory(nil)
	options := translateGetOptions(opts)
	result := factory.NewBuilder().
		Unstructured().
		NamespaceParam(options.Namespace).DefaultNamespace().AllNamespaces(options.AllNamespaces).
		FilenameParam(options.ExplicitNamespace, &options.FilenameOptions).
		LabelSelectorParam(options.LabelSelector).
		FieldSelectorParam(options.FieldSelector).
		ExportParam(options.Export).
		RequestChunksOf(options.ChunkSize).
		IncludeUninitialized(options.IncludeUninitialized).
		ResourceTypeOrNameArgs(true, typeOrName...).
		ContinueOnError().
		Latest().
		Flatten().
		Do()
	if err := result.Err(); err != nil {
		return "", err.Error()
	}
	object, err := result.Object()
	if err != nil {
		return "", err.Error()
	}
	payload, err := json.Marshal(object)
	if err != nil {
		return "", err.Error()
	}
	return string(payload), ""
}

//export Create
func Create(optsEncoded string) (res string, serr string) {
	defer func() {
		if r := recover(); r != nil {
			serr = fmt.Sprintf("%v\n%s", r, debug.Stack())
		}
	}()
	opts := map[string]interface{}{}
	if err := json.Unmarshal([]byte(optsEncoded), &opts); err != nil {
		return "", err.Error()
	}
	options := translateCreateOptions(opts)

	factory := cmdutil.NewFactory(nil)
	schema, err := factory.Validator(false)
	if err != nil {
		return "", err.Error()
	}

	cmdNamespace, enforceNamespace, err := factory.DefaultNamespace()
	if err != nil {
		return "", err.Error()
	}

	result := factory.NewBuilder().
		Unstructured().
		Schema(schema).
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &options.FilenameOptions).
		LabelSelectorParam(options.Selector).
		Flatten().
		Do()
	if err = result.Err(); err != nil {
		return "", err.Error()
	}

	object, err := result.Object()
	if err != nil {
		return "", err.Error()
	}
	payload, err := json.Marshal(object)
	if err != nil {
		return "", err.Error()
	}
	return string(payload), ""

	// count := 0
	// err = r.Visit(func(info *resource.Info, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err := kubectl.CreateOrUpdateAnnotation(
	// 		cmdutil.GetFlagBool(cmd, cmdutil.ApplyAnnotationsFlag),
	// 		info, f.JSONEncoder(),
	// 	)
	// 	if err != nil {
	// 		return cmdutil.AddSourceToErr("creating", info.Source, err)
	// 	}

	// 	if cmdutil.ShouldRecord(cmd, info) {
	// 		if err := cmdutil.RecordChangeCause(info.Object, f.Command(cmd, false)); err != nil {
	// 			return cmdutil.AddSourceToErr("creating", info.Source, err)
	// 		}
	// 	}

	// 	if !dryRun {
	// 		if err := createAndRefresh(info); err != nil {
	// 			return cmdutil.AddSourceToErr("creating", info.Source, err)
	// 		}
	// 	}

	// 	count++

	// 	shortOutput := output == "name"
	// 	if len(output) > 0 && !shortOutput {
	// 		return f.PrintResourceInfoForCommand(cmd, info, out)
	// 	}
	// 	if !shortOutput {
	// 		f.PrintObjectSpecificMessage(info.Object, out)
	// 	}

	// 	f.PrintSuccess(mapper, shortOutput, out, info.Mapping.Resource, info.Name, dryRun, "created")
	// 	return nil
	// })

	// if err != nil {
	// 	return "", err.Error()
	// }

	// payload, err := json.Marshal(object)
	// if err != nil {
	// 	return "", err.Error()
	// }
	// return string(payload), ""
}

func main() {}
