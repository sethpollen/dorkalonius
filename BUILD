load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_binary",
    "go_library",
    "go_test",
)

go_prefix("github.com/sethpollen/dorkalonius")

load("//tools:go_embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "coca_data",
    data = ["coca-5000.csv"],
    package = "dorkalonius",
)

go_library(
    name = "go_default_library",
    srcs = [
        "coca_word_list.go",
        "game.go",
        "memoize.go",
        "sampler.go",
        "sleep.go",
        "word_list.go",
        "word_set.go",
        ":coca_data",
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
    name = "word_list_test",
    srcs = ["word_list_test.go"],
    deps = [
        ":go_default_library",
    ],
)

go_test(
    name = "sampler_test",
    srcs = ["sampler_test.go"],
    deps = [
        ":go_default_library",
    ],
)

go_binary(
    name = "words_main",
    srcs = ["words_main.go"],
    deps = [
        ":go_default_library",
    ],
)
