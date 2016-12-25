load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "inflection_map.go",
        "inflection_xml.go",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "inflection_map_test",
    srcs = ["inflection_map_test.go"],
    deps = [":go_default_library"],
    data = ["inflections.xml.bz2"],
)