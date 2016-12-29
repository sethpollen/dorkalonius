load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")
load("//tools:go_embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "preference_data",
    data = ["inflection_to_base_preferences.csv"],
    package = "wiktionary",
)

go_library(
    name = "go_default_library",
    srcs = [
        "inflection_map.go",
        "inflection_xml.go",
        ":preference_data",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "inflection_map_test",
    srcs = ["inflection_map_test.go"],
    data = [":inflections_xml_bz2"],
    deps = [":go_default_library"],
)

filegroup(
    name = "inflections_xml_bz2",
    srcs = ["inflections.xml.bz2"],
    visibility = ["//visibility:public"],
)
