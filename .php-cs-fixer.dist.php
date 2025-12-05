<?php

$finder = (new PhpCsFixer\Finder())
    ->in([
        __DIR__ . '/src',
        __DIR__ . '/tests',
    ])
    ->exclude([
        'var',
        'vendor',
    ])
    ->notPath([
        'config/bundles.php',
        'config/preload.php',
    ])
;

return (new PhpCsFixer\Config())
    ->setRiskyAllowed(true)
    ->setRules([
        // Base rule sets
        '@Symfony' => true,
        '@Symfony:risky' => true,
        '@PHP82Migration' => true,
        '@PSR12' => true,

        // Strict types
        'declare_strict_types' => true,
        'strict_param' => true,
        'strict_comparison' => true,

        // Arrays
        'array_syntax' => ['syntax' => 'short'],
        'binary_operator_spaces' => [
            'default' => 'single_space',
        ],
        'no_multiline_whitespace_around_double_arrow' => true,
        'trim_array_spaces' => true,
        'whitespace_after_comma_in_array' => true,

        // Imports
        'ordered_imports' => [
            'imports_order' => ['class', 'function', 'const'],
            'sort_algorithm' => 'alpha',
        ],
        'global_namespace_import' => [
            'import_classes' => true,
            'import_constants' => false,
            'import_functions' => false,
        ],
        'no_unused_imports' => true,

        // Classes and methods
        'class_attributes_separation' => [
            'elements' => [
                'const' => 'one',
                'method' => 'one',
                'property' => 'one',
                'trait_import' => 'none',
            ],
        ],
        'final_class' => true,
        'final_public_method_for_abstract_class' => true,
        'no_null_property_initialization' => true,
        'self_static_accessor' => true,

        // Control structures
        'no_superfluous_elseif' => true,
        'no_useless_else' => true,
        'simplified_if_return' => true,

        // Functions
        'native_function_invocation' => false,
        'void_return' => true,

        // PHPDoc
        'phpdoc_align' => ['align' => 'left'],
        'phpdoc_order' => true,
        'phpdoc_separation' => true,
        'phpdoc_summary' => true,
        'phpdoc_trim' => true,
        'phpdoc_types_order' => [
            'null_adjustment' => 'always_last',
            'sort_algorithm' => 'alpha',
        ],

        // Comments
        'no_empty_comment' => true,
        'single_line_comment_style' => true,

        // Code quality
        'modernize_types_casting' => true,
        'no_useless_return' => true,
        'combine_consecutive_issets' => true,
        'combine_consecutive_unsets' => true,
        'is_null' => true,
        'logical_operators' => true,
        'no_alias_functions' => true,

        // Whitespace and formatting
        'blank_line_before_statement' => [
            'statements' => ['return', 'throw', 'try', 'if', 'for', 'foreach', 'while', 'switch'],
        ],
        'yoda_style' => false,
        'concat_space' => ['spacing' => 'one'],
        'method_chaining_indentation' => true,
        'no_extra_blank_lines' => [
            'tokens' => [
                'extra',
                'throw',
                'use',
            ],
        ],
        'no_spaces_around_offset' => true,
        'types_spaces' => ['space' => 'none'],
    ])
    ->setFinder($finder)
    ->setCacheFile(__DIR__ . '/var/.php-cs-fixer.cache')
;
