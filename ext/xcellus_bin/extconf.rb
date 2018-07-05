require 'mkmf'

extension_name = 'xcellus_bin'
target_platform = 'linux'
if RbConfig::CONFIG['target'] =~ /darwin/
  target_platform = 'darwin'
end

bindings = "xcbindings-#{target_platform}-amd64"

dir_config(extension_name)
dir_config(bindings, '$(srcdir)/../../bin', '$(srcdir)/../../bin')
$LDFLAGS << " -l#{bindings} -Wl,-rpath,./../../bin -Wl,-rpath,./../bin"
create_makefile(extension_name)
