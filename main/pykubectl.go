package main

import (
	"encoding/json"

	// You should never import this, this will try to reinitialize
	// flags, and panic
	//
	// utilflag "k8s.io/apiserver/pkg/util/flag"

	"k8s.io/kubernetes/pkg/kubectl/cmd/resource"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	kresource "k8s.io/kubernetes/pkg/kubectl/resource"
)
import "C"

func translateOptions(raw map[string]interface{}) *resource.GetOptions {
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

//export ResourceGet
func ResourceGet(optsEncoded string, typeOrName []string) (string, string) {
	opts := map[string]interface{}{}
	if err := json.Unmarshal([]byte(optsEncoded), &opts); err != nil {
		return "", err.Error()
	}
	factory := cmdutil.NewFactory(nil)
	options := translateOptions(opts)
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

func main() {}
