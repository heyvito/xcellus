require File.expand_path('../xcellus_bin', __FILE__)

require 'json'

module Xcellus
  VERSION = '0.1.2'.freeze

  class << self
    def transform(input)
      unless input.kind_of? Array
        raise ArgumentError, 'Xcellus.transform only accepts Arrays'
      end
      StringIO.new(Xcellus::process(input.to_json))
    end
  end
end
