load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_binary",
    "go_library",
    "go_test",
)

go_prefix("github.com/sethpollen/dorkalonius")

go_library(
    name = "go_default_library",
    srcs = [
        "sampler.go",
        "sleep.go",
        "word_list.go",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "word_list_test",
    srcs = ["word_list_test.go"],
    deps = [
        ":go_default_library",
        "//coca:go_default_library",
    ],
)

go_test(
    name = "sampler_test",
    srcs = ["sampler_test.go"],
    deps = [
        ":go_default_library",
        "//coca:go_default_library",
    ],
)

go_binary(
    name = "words_main",
    srcs = ["words_main.go"],
    deps = [
        ":go_default_library",
        "//coca:go_default_library",
    ],
)
