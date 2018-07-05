#include <stdbool.h>
#include <ruby.h>
#include "../../bin/libxcbindings.h"

VALUE rb_process(VALUE self, VALUE str) {
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

void Init_xcellus_bin() {
    VALUE module = rb_define_module("Xcellus");
    rb_define_singleton_method(module, "process", rb_process, 1);
}
