#include <stdbool.h>
#include <ruby.h>
#include "../../bin/libxcbindings.h"

VALUE rb_xcellus_transform(VALUE self, VALUE str) {
	if (RB_TYPE_P(str, T_STRING) == 0) return Qnil;
	struct go_xcellus_process_return result = go_xcellus_process(RSTRING_PTR(str));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	GoInt buffer_len = result.r2;
	char *buffer = result.r3;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}

    VALUE retval = rb_str_new(buffer, buffer_len);
    free(buffer);
    return retval;
}

VALUE rb_xcellus_load(VALUE self, VALUE path) {
	if (RB_TYPE_P(path, T_STRING) == 0) return Qnil;
	struct go_xcellus_load_return result = go_xcellus_load(RSTRING_PTR(path));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	char *handle = result.r2;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	VALUE retval = rb_str_new_cstr(handle);
	return retval;
}

VALUE rb_xcellus_find_in_column(VALUE self, VALUE handle, VALUE sheetName, VALUE index, VALUE value) {
	if (RB_TYPE_P(handle, T_STRING) == 0) return Qnil;
	if (RB_TYPE_P(sheetName, T_STRING) == 0) return Qnil;
	if (RB_TYPE_P(index, T_FIXNUM) == 0) return Qnil;
	if (RB_TYPE_P(value, T_STRING) == 0) return Qnil;
	struct go_xcellus_find_in_column_return result = go_xcellus_find_in_column(RSTRING_PTR(handle), RSTRING_PTR(sheetName), RSTRING_PTR(value), FIX2INT(index));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	long long foundIndex = result.r2;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	return INT2NUM(foundIndex);
}

VALUE rb_xcellus_replace_row(VALUE self, VALUE handle, VALUE sheetName, VALUE data, VALUE index) {
	if (RB_TYPE_P(handle, T_STRING) == 0) return Qnil;
	if (RB_TYPE_P(sheetName, T_STRING) == 0) return Qnil;
	if (RB_TYPE_P(index, T_FIXNUM) == 0) return Qnil;
	if (RB_TYPE_P(data, T_STRING) == 0) return Qnil;
	struct go_xcellus_replace_row_return result = go_xcellus_replace_row(RSTRING_PTR(handle), RSTRING_PTR(sheetName), RSTRING_PTR(data), FIX2INT(index));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	return Qtrue;
}

VALUE rb_xcellus_close(VALUE self, VALUE handle) {
	if (RB_TYPE_P(handle, T_STRING) == 0) return Qnil;
	struct go_xcellus_end_return result = go_xcellus_end(RSTRING_PTR(handle));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	return Qtrue;
}

VALUE rb_xcellus_save(VALUE self, VALUE handle, VALUE path) {
	if (RB_TYPE_P(handle, T_STRING) == 0) return Qnil;
	if (RB_TYPE_P(path, T_STRING) == 0) return Qnil;
	struct go_xcellus_save_return result = go_xcellus_save(RSTRING_PTR(handle), RSTRING_PTR(path));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	return Qtrue;
}

VALUE rb_xcellus_append(VALUE self, VALUE handle, VALUE data) {
	if (RB_TYPE_P(handle, T_STRING) == 0 || RB_TYPE_P(data, T_STRING) == 0)
		return Qnil;
	struct go_xcellus_append_return result = go_xcellus_append(RSTRING_PTR(handle), RSTRING_PTR(data));
	bool failed = result.r0 == 1;
	char *error_description = result.r1;
	if(failed) {
		VALUE errObj = rb_exc_new2(rb_eStandardError, error_description);
		free(error_description);
		rb_exc_raise(errObj);
	}
	return Qtrue;
}

void Init_xcellus_bin() {
	VALUE module = rb_define_module("Xcellus");

	rb_define_singleton_method(module, "_transform", rb_xcellus_transform, 1);
	rb_define_singleton_method(module, "_load", rb_xcellus_load, 1);

	rb_define_singleton_method(module, "_find_in_column", rb_xcellus_find_in_column, 4);
	rb_define_singleton_method(module, "_replace_row", rb_xcellus_replace_row, 4);
	rb_define_singleton_method(module, "_save", rb_xcellus_save, 2);
	rb_define_singleton_method(module, "_close", rb_xcellus_close, 2);
	rb_define_singleton_method(module, "_append", rb_xcellus_append, 2);
}
