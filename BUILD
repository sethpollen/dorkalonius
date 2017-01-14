load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_binary",
    "go_library",
    "go_test",
)

go_prefix("github.com/sethpollen/dorkalonius")

load("//tools:tools.bzl", "go_embed_data", "word_set")

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
    name = "go_default_library",
    srcs = [
        "game.go",
        ":coca_adjective_set",
        ":coca_word_set",
    ],
    deps = [
        "//util:go_default_library",
    ],
)

go_binary(
    name = "game_test_main",
    srcs = ["game_test_main.go"],
    deps = [
        "//util:go_default_library",
        ":go_default_library",
    ],
)

go_binary(
    name = "words_main",
    srcs = ["words_main.go"],
    deps = [
        "//util:go_default_library",
        ":go_default_library",
    ],
)
