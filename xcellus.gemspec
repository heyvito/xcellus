Gem::Specification.new do |s|
  s.name    = 'xcellus'
  s.version = '2.0'
  s.summary       = %q{Xcellus is an ugly XLSX generator}
  s.description   = <<-description
    Xcellus leverages a native extension to create XLSX files. It is an 'ugly'
    generator because it lacks all styling mechanisms to make your XLSX pretty.
  description
  s.authors  = ['Victor Gama']
  s.email    = ['hey@vito.io']
  s.homepage      = 'https://github.com/victorgama/xcellus'
  s.license       = 'MIT'

  s.files = Dir.glob('ext/**/*.{c,rb}') +
            Dir.glob('ext/Rakefile') +
            Dir.glob('lib/**/*.rb') +
            Dir.glob('bin/*.{h,so}')

  s.extensions << 'ext/xcellus_bin/extconf.rb'
  s.extensions << 'ext/Rakefile'

  s.add_development_dependency 'rake-compiler'
  s.add_development_dependency 'bundler', '~> 1.16'
  s.add_development_dependency 'rake', '~> 12.0'
end
