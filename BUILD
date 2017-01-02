load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_binary",
    "go_library",
    "go_test",
)

go_prefix("github.com/sethpollen/dorkalonius")

load("//tools:tools.bzl", "go_embed_data", "word_set")

# TODO: remove and use word_set instead
go_embed_data(
    name = "coca_data",
    data = ["coca-5000.csv"],
    package = "dorkalonius",
)

word_set(
    name = "coca_data2",
    srcs = ["coca-5000.csv"],
    csv_header_lines = 2,
    csv_word_column = 1,
    csv_weight_column = 3,
)

go_library(
    name = "go_default_library",
    srcs = [
        "coca_word_list.go",
        "game.go",
        "memoize.go",
        "sleep.go",
        "word_set.go",
        "word_set_builder.go",
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
    name = "word_set_test",
    srcs = ["word_set_test.go"],
    deps = [
        ":go_default_library",
    ],
)

go_binary(
    name = "game_test_main",
    srcs = ["game_test_main.go"],
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
