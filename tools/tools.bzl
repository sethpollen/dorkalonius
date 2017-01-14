# Bazel rule for embedding data into a Go program.
#
# Arguments:
#   name: Name of the rule.
#   data: List of targets providing the files to embed.
#   package: Go package to put the generated code in.
def go_embed_data(
  name,
  data,
  package
):
  native.genrule(
    name = name,
    srcs = data,
    outs = [name + "_embed_data.go"],
    cmd = "./$(location //tools:go_embed_main)" +
          "  --output_file=\"$@\"" +
          "  --package=\"" + package + "\"" +
          "  --method=\"Get_" + name + "\"" +
          "  $(SRCS)",
    tools = ["//tools:go_embed_main"],
  )
  
# Bazel rule for converting a CSV file of word count data into a set of .go file
# which provide a WordSet object.
#
# Arguments:
#   name: Name of the rule
#   srcs: List of targets providing input CSV files
#   package: Go package to put the generated code in.
def word_set(
  name,
  srcs,
  package,
  csv_header_lines=0,
  csv_word_column=0,
  csv_weight_column=1,
  csv_filter_column=None,
  csv_filter_value=None,
):
  filter_flags = ""
  if csv_filter_column and csv_filter_value:
    filter_flags = "--csv_filter_column=%d --csv_filter_value=%s" % (
        csv_filter_column, csv_filter_value)
  native.genrule(
    name = name + "__wordset",
    srcs = srcs,
    outs = [name + ".wordset"],
    cmd = "./$(location //tools:csv_to_word_set_main)" +
          "  --output_file=\"$@\"" +
          "  --csv_header_lines=" + str(csv_header_lines) +
          "  --csv_word_column=" + str(csv_word_column) +
          "  --csv_weight_column=" + str(csv_weight_column) +
          "  " + filter_flags +
          "  $(SRCS)",
    tools = ["//tools:csv_to_word_set_main"],
  )
  go_embed_data(
    name = name + "__embed",
    data = [":" + name + "__wordset"],
    package = package,
  )
  # We emit another .go file to handle deserializing and memoizing the WordSet.
  go_lib_source = " ; ".join([
    'package %s' % package,
    'import "github.com/sethpollen/dorkalonius/util"',
    'import "log"',
    'var %s_memo = util.NewMemo(func() interface{} {' % name,
    'in := Get_%s__embed()' % name,
    'wordSet, err := util.DeserializeWordSet(in)',
    'if err != nil {',
    'log.Fatal(err)',
    '}',
    'return wordSet',
    '})',
    'func Get_%s() *util.WordSet {' % name,
    'return %s_memo.Get().(*util.WordSet)' % name,
    '}',
  ])
  native.genrule(
    name = name + "__library_source",
    outs = [name + ".go"],
    cmd = "echo '%s' > \"$@\"" % go_lib_source,
  )
  native.filegroup(
    name = name,
    srcs = [
      ":" + name + "__embed",
      ":" + name + "__library_source",
    ],
  )
