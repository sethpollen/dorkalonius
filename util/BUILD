load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "memoize.go",
        "sleep.go",
        "word_set.go",
        "word_set_builder.go",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "memoize_test",
    srcs = ["memoize_test.go"],
    deps = [
        ":go_default_library",
    ],
)

go_test(
    name = "word_set_test",
    srcs = ["word_set_test.go"],
    deps = [
        ":go_default_library",
    ],
)