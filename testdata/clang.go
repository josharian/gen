package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

/**
 * \brief Return the CXTranslationUnit for a given source file and the provided
 * command line arguments one would pass to the compiler.
 *
 * Note: The 'source_filename' argument is optional.  If the caller provides a
 * NULL pointer, the name of the source file is expected to reside in the
 * specified command line arguments.
 *
 * Note: When encountered in 'clang_command_line_args', the following options
 * are ignored:
 *
 *   '-c'
 *   '-emit-ast'
 *   '-fsyntax-only'
 *   '-o <output file>'  (both '-o' and '<output file>' are ignored)
 *
 * \param CIdx The index object with which the translation unit will be
 * associated.
 *
 * \param source_filename - The name of the source file to load, or NULL if the
 * source file is included in \p clang_command_line_args.
 *
 * \param num_clang_command_line_args The number of command-line arguments in
 * \p clang_command_line_args.
 *
 * \param clang_command_line_args The command-line arguments that would be
 * passed to the \c clang executable if it were being invoked out-of-process.
 * These command-line options will be parsed and will affect how the translation
 * unit is parsed. Note that the following options are ignored: '-c',
 * '-emit-ast', '-fsyntex-only' (which is the default), and '-o <output file>'.
 *
 * \param num_unsaved_files the number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param unsaved_files the files that have not yet been saved to disk
 * but may be required for code completion, including the contents of
 * those files.  The contents and name of these files (as specified by
 * CXUnsavedFile) are copied when necessary, so the client only needs to
 * guarantee their validity until the call to this function returns.
 */
func (idx Index) CreateTranslationUnitFromSourceFile(fname string, args []string, us UnsavedFiles) TranslationUnit {
	var (
		c_fname *C.char = nil
		c_us            = us.to_c()
	)
	defer c_us.Dispose()
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i, _ := range args {
		cstr := C.CString(args[i])
		defer C.free(unsafe.Pointer(cstr))
		c_cmds[i] = cstr
	}

	var c_argv **C.char = nil
	if c_nargs > 0 {
		c_argv = &c_cmds[0]
	}

	o := C.clang_createTranslationUnitFromSourceFile(
		idx.c,
		c_fname,
		c_nargs, c_argv,
		C.uint(len(c_us)), c_us.ptr())
	return TranslationUnit{o}

}

/**
 * \brief Parse the given source file and the translation unit corresponding
 * to that file.
 *
 * This routine is the main entry point for the Clang C API, providing the
 * ability to parse a source file into a translation unit that can then be
 * queried by other functions in the API. This routine accepts a set of
 * command-line arguments so that the compilation can be configured in the same
 * way that the compiler is configured on the command line.
 *
 * \param CIdx The index object with which the translation unit will be
 * associated.
 *
 * \param source_filename The name of the source file to load, or NULL if the
 * source file is included in \p command_line_args.
 *
 * \param command_line_args The command-line arguments that would be
 * passed to the \c clang executable if it were being invoked out-of-process.
 * These command-line options will be parsed and will affect how the translation
 * unit is parsed. Note that the following options are ignored: '-c',
 * '-emit-ast', '-fsyntex-only' (which is the default), and '-o <output file>'.
 *
 * \param num_command_line_args The number of command-line arguments in
 * \p command_line_args.
 *
 * \param unsaved_files the files that have not yet been saved to disk
 * but may be required for parsing, including the contents of
 * those files.  The contents and name of these files (as specified by
 * CXUnsavedFile) are copied when necessary, so the client only needs to
 * guarantee their validity until the call to this function returns.
 *
 * \param num_unsaved_files the number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param options A bitmask of options that affects how the translation unit
 * is managed but not its compilation. This should be a bitwise OR of the
 * CXTranslationUnit_XXX flags.
 *
 * \returns A new translation unit describing the parsed code and containing
 * any diagnostics produced by the compiler. If there is a failure from which
 * the compiler cannot recover, returns NULL.
 */
func (idx Index) Parse(fname string, args []string, us UnsavedFiles, options TranslationUnitFlags) TranslationUnit {
	var (
		c_fname *C.char = nil
		c_us            = us.to_c()
	)
	defer c_us.Dispose()
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i, _ := range args {
		cstr := C.CString(args[i])
		defer C.free(unsafe.Pointer(cstr))
		c_cmds[i] = cstr
	}

	var c_args **C.char = nil
	if len(args) > 0 {
		c_args = &c_cmds[0]
	}
	o := C.clang_parseTranslationUnit(
		idx.c,
		c_fname,
		c_args, c_nargs,
		c_us.ptr(), C.uint(len(c_us)),
		C.uint(options))
	return TranslationUnit{o}

}

// EOF
