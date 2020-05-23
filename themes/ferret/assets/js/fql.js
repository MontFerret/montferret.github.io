/*
Language: FQL
Category: common, scripting
*/
function registerFQL(hljs) {
    var IDENT_RE = '[A-Za-z$_][0-9A-Za-z$_]*';
    var KEYWORDS = {
      keyword:
        'USE FOR IN RETURN LET AND OR LIMIT FILTER DISTINCT SORT COLLECT ASC DESC' +
        'INTO KEEP WITH COUNT ALL ANY AGGREGATE LIKE NOT'
      ,
      literal:
        'TRUE true FALSE false NONE',
    };
    var NUMBER = {
      className: 'number',
      variants: [
        { begin: '\\b(0[bB][01]+)' },
        { begin: '\\b(0[oO][0-7]+)' },
        { begin: hljs.C_NUMBER_RE }
      ],
      relevance: 0
    };
    var SUBST = {
      className: 'subst',
      begin: '\\$\\{', end: '\\}',
      keywords: KEYWORDS,
      contains: []  // defined later
    };
    var TEMPLATE_STRING = {
      className: 'string',
      begin: '`', end: '`',
      contains: [
        hljs.BACKSLASH_ESCAPE,
        SUBST
      ]
    };
    SUBST.contains = [
      hljs.APOS_STRING_MODE,
      hljs.QUOTE_STRING_MODE,
      TEMPLATE_STRING,
      NUMBER,
      hljs.REGEXP_MODE
    ]
    var PARAMS_CONTAINS = SUBST.contains.concat([
      hljs.C_BLOCK_COMMENT_MODE,
      hljs.C_LINE_COMMENT_MODE
    ]);
  
    return {
      aliases: ['fql'],
      case_insensitive: true,
      keywords: KEYWORDS,
      contains: [
        hljs.APOS_STRING_MODE,
        hljs.QUOTE_STRING_MODE,
        TEMPLATE_STRING,
        hljs.C_LINE_COMMENT_MODE,
        hljs.C_BLOCK_COMMENT_MODE,
        NUMBER,
        { // object attr container
          begin: /[{,]\s*/, relevance: 0,
          contains: [
            {
              begin: IDENT_RE + '\\s*:', returnBegin: true,
              relevance: 0,
              contains: [{className: 'attr', begin: IDENT_RE, relevance: 0}]
            }
          ]
        },
        { // "value" container
          begin: '(' + hljs.RE_STARTERS_RE + '|\\b(case|return|throw)\\b)\\s*',
          keywords: 'return throw case',
          contains: [
            hljs.C_LINE_COMMENT_MODE,
            hljs.C_BLOCK_COMMENT_MODE,
            hljs.REGEXP_MODE,
            {
              className: 'function',
              begin: '(\\(.*?\\)|' + IDENT_RE + ')\\s*=>', returnBegin: true,
              end: '\\s*=>',
              contains: [
                {
                  className: 'params',
                  variants: [
                    {
                      begin: IDENT_RE
                    },
                    {
                      begin: /\(\s*\)/,
                    },
                    {
                      begin: /\(/, end: /\)/,
                      excludeBegin: true, excludeEnd: true,
                      keywords: KEYWORDS,
                      contains: PARAMS_CONTAINS
                    }
                  ]
                }
              ]
            }
          ],
          relevance: 0
        },
        hljs.METHOD_GUARD,
      ],
      illegal: /#(?!!)/
    };
  }