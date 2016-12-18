load("@io_bazel_rules_go//go:def.bzl", "go_library")

filegroup(
    name = "data",
    srcs = ["inflections.xml"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "inflection_xml.go",
    ],
    data = [":data"],
    visibility = ["//visibility:public"],
)
