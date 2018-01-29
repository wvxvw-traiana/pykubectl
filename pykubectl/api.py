# -*- coding: utf-8 -*-
import json

from pykubectl.wrapped import pykubectl_get_impl, pykubectl_create_impl


def kw_to_dict(f, locs):
    args = f.__code__.co_varnames[:f.__code__.co_argcount]
    return {k: v for k, v in locs.items() if k in args and v is not None}


def kubectl_get(
        items,
        filenames=None,
        recursive=None,
        watch=None,
        raw=None,
        watch_only=None,
        chunk_size=None,
        label_selector=None,
        field_selector=None,
        all_namespaces=None,
        namespace=None,
        explicit_namespace=None,
        ignore_not_found=None,
        show_kind=None,
        export=None,
        include_uninitialized=None,
):
    options = kw_to_dict(kubectl_get, locals())
    options = json.dumps(options).encode("utf-8")
    return json.loads(pykubectl_get_impl(items, options))


def kubectl_create(
        filenames=None,
        recursive=None,
        raw=None,
        edit_before_create=None,
        selector=None,
):
    options = kw_to_dict(kubectl_create, locals())
    options = json.dumps(options).encode("utf-8")
    return json.loads(pykubectl_create_impl(options))
