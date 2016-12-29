load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("//tools:go_embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "coca_data",
    data = ["coca-5000.csv"],
    package = "coca",
)

go_library(
    name = "go_default_library",
    srcs = [
        "word_list.go",
        ":coca_data",
    ],
    visibility = ["//visibility:public"],
    deps = ["//:go_default_library"],
)
