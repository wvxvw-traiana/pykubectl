package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	// You should never import this, this will try to reinitialize
	// flags, and panic
	//
	// utilflag "k8s.io/apiserver/pkg/util/flag"

	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/resource"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	kresource "k8s.io/kubernetes/pkg/kubectl/resource"
)
import "C"

func errWithStack(original error) *C.char {
	scanner := bufio.NewScanner(bytes.NewReader(debug.Stack()))
	i := 0
	lines := []string{original.Error()}
	for scanner.Scan() {
		if i > 4 {
			lines = append(lines, scanner.Text())
		}
		i++
	}
	return C.CString(strings.Join(lines, "\n"))
}

func translateFilenames(raw map[string]interface{}) kresource.FilenameOptions {
	foptions := kresource.FilenameOptions{}

	if fnames, ok := raw["filenames"]; ok {
		rawNames := fnames.([]interface{})
		names := make([]string, len(rawNames))
		for i := 0; i < len(names); i++ {
			names[i] = rawNames[i].(string)
		}
		foptions.Filenames = names
	}
	if recursive, ok := raw["recursive"]; ok {
		foptions.Recursive = recursive.(bool)
	}
	return foptions
}

func translateGetOptions(raw map[string]interface{}) *resource.GetOptions {

	result := &resource.GetOptions{
		FilenameOptions: translateFilenames(raw),
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
	result := &cmd.CreateOptions{
		FilenameOptions: translateFilenames(raw),
	}
	for k, v := range raw {
		switch k {
		case "raw":
			result.Raw = v.(string)
		case "edit_before_create":
			result.EditBeforeCreate = v.(bool)
		case "selector":
			result.Selector = v.(string)
		}
	}
	return result
}

//export ResourceGet
func ResourceGet(optsEncoded string, typeOrName []string) (res *C.char, serr *C.char) {
	defer func() {
		if r := recover(); r != nil {
			serr = C.CString(fmt.Sprintf("%v\n%s", r, debug.Stack()))
		}
	}()
	empty := C.CString("")
	opts := map[string]interface{}{}
	fmt.Printf("Decoding: %v\n", optsEncoded)
	fmt.Printf("Type or name: %v\n", typeOrName)
	if err := json.Unmarshal([]byte(optsEncoded), &opts); err != nil {
		return empty, errWithStack(err)
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
		return empty, errWithStack(err)
	}
	object, err := result.Object()
	if err != nil {
		return empty, errWithStack(err)
	}
	payload, err := json.Marshal(object)
	if err != nil {
		return empty, errWithStack(err)
	}
	return C.CString(string(payload)), empty
}

func createAndRefresh(info *kresource.Info) error {
	obj, err := kresource.NewHelper(info.Client, info.Mapping).
		Create(info.Namespace, true, info.Object)
	if err != nil {
		return err
	}
	info.Refresh(obj, true)
	return nil
}

//export Create
func Create(optsEncoded string) (res *C.char, serr *C.char) {
	defer func() {
		if r := recover(); r != nil {
			serr = C.CString(fmt.Sprintf("%v\n%s", r, debug.Stack()))
		}
	}()
	empty := C.CString("")
	opts := map[string]interface{}{}
	if err := json.Unmarshal([]byte(optsEncoded), &opts); err != nil {
		return empty, errWithStack(err)
	}
	options := translateCreateOptions(opts)

	factory := cmdutil.NewFactory(nil)
	schema, err := factory.Validator(false)
	if err != nil {
		return empty, errWithStack(err)
	}

	cmdNamespace, enforceNamespace, err := factory.DefaultNamespace()
	if err != nil {
		return empty, errWithStack(err)
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
		return empty, errWithStack(err)
	}

	// TODO(olegs): This doesn't come from otpions struct, rather from
	// global command-line arguments, need to pass these somehow too.
	dryRun := false

	count := 0
	err = result.Visit(func(info *kresource.Info, err error) error {
		if err != nil {
			return err
		}
		// TODO(olegs): I don't really know what this "true" means
		e := kubectl.CreateOrUpdateAnnotation(true, info, factory.JSONEncoder())
		if e != nil {
			return cmdutil.AddSourceToErr("creating", info.Source, err)
		}

		if !dryRun {
			if e := createAndRefresh(info); e != nil {
				return cmdutil.AddSourceToErr("creating", info.Source, e)
			}
		}

		count++
		return nil
	})

	if err != nil {
		return empty, errWithStack(err)
	}

	object, err := result.Object()
	if err != nil {
		return empty, errWithStack(err)
	}

	payload, err := json.Marshal(object)
	if err != nil {
		return empty, errWithStack(err)
	}
	return C.CString(string(payload)), empty
}

func main() {}
