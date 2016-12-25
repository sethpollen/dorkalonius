load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")
load("//tools:go_embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "inflections_xml",
    data = ["inflections.xml"],
    package = "embed",
)

go_library(
    name = "go_default_library",
    srcs = [":inflections_xml"],
    visibility = ["//visibility:public"],
    deps = ["//incubator/wiktionary:go_default_library"],
)

go_test(
    name = "embed_test",
    srcs = ["embed_test.go"],
    deps = [":go_default_library"],
)
