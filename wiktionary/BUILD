load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")
load("//tools:go_embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "preferences_csv",
    data = ["preferences.csv"],
    package = "wiktionary",
)

go_library(
    name = "go_default_library",
    srcs = [
        "inflection_map.go",
        "inflection_xml.go",
        ":preferences_csv",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "inflection_map_test",
    srcs = ["inflection_map_test.go"],
    data = ["inflections.xml.bz2"],
    deps = [":go_default_library"],
)
