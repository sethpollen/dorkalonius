load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_binary",
    "go_library",
    "go_test",
)

go_prefix("github.com/sethpollen/dorkalonius")

load("//tools:tools.bzl", "go_embed_data", "word_set")

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

word_set(
    name = "coca_word_set",
    srcs = ["coca-5000.csv"],
    package = "dorkalonius",
    csv_header_lines = 2,
    csv_word_column = 1,
    csv_weight_column = 3,
)

word_set(
    name = "coca_adjective_set",
    srcs = ["coca-5000.csv"],
    package = "dorkalonius",
    csv_header_lines = 2,
    csv_word_column = 1,
    csv_weight_column = 3,
    csv_filter_column = 2,
    csv_filter_value = "j"
)

go_library(
    name = "game",
    srcs = [
        ":coca_adjective_set",
        ":coca_word_set",
        "game.go",
    ],
    deps = [
        ":go_default_library",
    ],
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
        ":game",
    ],
)

go_binary(
    name = "words_main",
    srcs = ["words_main.go"],
    deps = [
        ":go_default_library",
        ":game",
    ],
)
