ace.define(
  "ace/mode/fql_highlight_rules",
  [
    "require",
    "exports",
    "module",
    "ace/lib/oop",
    "ace/mode/text_highlight_rules"
  ],
  function(require, exports, module) {
    "use strict";

    var oop = require("../lib/oop");
    var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;

    var FqlHighlightRules = function() {
      var identifier = "[A-Za-z_][0-9A-Za-z_]*";
      var keywords = [
        "USE", "AS", "MATCH", "WHEN", "FUNC", "FOR", "RETURN", "QUERY",
        "USING", "WAITFOR", "DISPATCH", "OPTIONS", "TIMEOUT", "EVERY",
        "BACKOFF", "JITTER", "EXISTS", "COUNT", "VALUE", "ONE",
        "DISTINCT", "FILTER", "SORT", "LIMIT", "LET", "VAR", "COLLECT",
        "ASC", "DESC", "AT", "LEAST", "INTO", "KEEP", "WITH", "ALL",
        "ANY", "AGGREGATE", "EVENT", "LIKE", "NOT", "IN", "DO", "WHILE",
        "AND", "OR", "ON", "ERROR", "FAIL", "RETRY", "DELAY"
      ].join("|");
      var literals = [
        "TRUE", "FALSE", "NONE", "NULL"
      ].join("|");
      var keywordMapper = this.createKeywordMapper({
        "keyword": keywords,
        "constant.language": literals
      }, "identifier", true);

      this.$rules = {
        start: [
          {
            token: "comment",
            regex: "\\/\\/.*$"
          },
          {
            token: "comment.start",
            regex: "\\/\\*",
            next: "comment"
          },
          {
            token: "string",
            regex: '"(?:\\\\.|""|[^"\\\\])*"'
          },
          {
            token: "string",
            regex: "'(?:\\\\.|''|[^'\\\\])*'"
          },
          {
            token: "string",
            regex: "`(?:\\\\.|\\$\\{[^}]*\\}|[^`\\\\])*`"
          },
          {
            token: "string",
            regex: "\\xB4(?:\\\\\\xB4|[^\\xB4])*\\xB4"
          },
          {
            token: "variable.parameter",
            regex: "@" + identifier + "\\b"
          },
          {
            token: "constant.numeric",
            regex: "\\b[0-9]+(?:\\.[0-9]+)?(?:[eE][+-]?[0-9]+)?(?:MS|S|M|H|D)\\b"
          },
          {
            token: "constant.numeric",
            regex: "\\b(?:[0-9]+\\.[0-9]+(?:[eE][+-]?[0-9]+)?|[0-9]+(?:[eE][+-]?[0-9]+)?)\\b"
          },
          {
            token: "support.type",
            regex: identifier + "(?=::)"
          },
          {
            token: "keyword.operator",
            regex: "::|\\.\\.|\\+\\+|--|\\*=|\\/=|\\+=|-=|==|!=|>=|<=|=~|!~|<-|=>|&&|\\|\\||[+\\-*/%<>~=!?]"
          },
          {
            token: "paren.lparen",
            regex: "[\\[{(]"
          },
          {
            token: "paren.rparen",
            regex: "[\\]})]"
          },
          {
            token: "punctuation.operator",
            regex: "[:,.;]"
          },
          {
            token: function(value) {
              var token = keywordMapper(value);

              if (token === "identifier") {
                return "support.function";
              }

              return token;
            },
            regex: identifier + "(?=\\s*\\()"
          },
          {
            token: keywordMapper,
            regex: identifier + "\\b"
          },
          {
            token: "text",
            regex: "\\s+"
          }
        ],
        comment: [
          {
            token: "comment.end",
            regex: "\\*\\/",
            next: "start"
          },
          {
            defaultToken: "comment"
          }
        ]
      };

      this.normalizeRules();
    };

    oop.inherits(FqlHighlightRules, TextHighlightRules);

    exports.FqlHighlightRules = FqlHighlightRules;
  }
);

ace.define(
  "ace/mode/fql",
  [
    "require",
    "exports",
    "module",
    "ace/lib/oop",
    "ace/mode/text",
    "ace/mode/fql_highlight_rules"
  ],
  function(require, exports, module) {
    "use strict";

    var oop = require("../lib/oop");
    var TextMode = require("./text").Mode;
    var FqlHighlightRules = require("./fql_highlight_rules").FqlHighlightRules;

    var Mode = function() {
      this.HighlightRules = FqlHighlightRules;
      this.$behaviour = this.$defaultBehaviour;
    };

    oop.inherits(Mode, TextMode);

    (function() {
      this.lineCommentStart = "//";
      this.blockComment = {
        start: "/*",
        end: "*/"
      };
      this.$id = "ace/mode/fql";
    }).call(Mode.prototype);

    exports.Mode = Mode;
  }
);

ace.define(
  "ace/mode/aql_highlight_rules",
  [
    "require",
    "exports",
    "module",
    "ace/mode/fql_highlight_rules"
  ],
  function(require, exports, module) {
    "use strict";

    exports.AqlHighlightRules = require("./fql_highlight_rules").FqlHighlightRules;
  }
);

ace.define(
  "ace/mode/aql",
  [
    "require",
    "exports",
    "module",
    "ace/lib/oop",
    "ace/mode/fql"
  ],
  function(require, exports, module) {
    "use strict";

    var oop = require("../lib/oop");
    var FqlMode = require("./fql").Mode;

    var Mode = function() {
      FqlMode.call(this);
      this.$id = "ace/mode/aql";
    };

    oop.inherits(Mode, FqlMode);

    exports.Mode = Mode;
  }
);

(function() {
  ace.require(["ace/mode/fql"], function(mode) {
    if (typeof module === "object" && typeof exports === "object" && module) {
      module.exports = mode;
    }
  });
})();
