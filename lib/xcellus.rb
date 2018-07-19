require File.expand_path('../xcellus_bin', __FILE__)

require 'json'

# Xcellus provides a clean interface to an underlying, native XLSX parser/writer
module Xcellus
  VERSION = '2.0'.freeze

  class << self
    # Transforms a provided array of objects into a XLSX file, and returns it as
    # a StringIO object.
    # Example:
    # Xcellus.transform [
    #   {
    #     title: 'Brian\'s Worksheet',
    #     headers: [:Artist, :Track, :Playcount],
    #     rows: [
    #         ['Metallica',               'Hero of the Day',      242],
    #         ['Metallica',               'The Shortest Straw',   186],
    #         ['Queens of the Stone Age', 'My God Is the Sun',    276],
    #         ['Queens of the Stone Age', 'I Sat by the Ocean',   270],
    #         ['Gorillaz',                'On Melancholy Hill',   203],
    #         ['Gorillaz',                'Kids With Guns',       184],
    #     ]
    #   }
    # ]
    # => StringIO
    # The returned XLSX file will contain a single worksheet named "Brian's
    # Worksheet", and the provided headers and rows.
    def transform(input)
      unless input.kind_of? Array
        raise ArgumentError, 'Xcellus.transform only accepts Arrays'
      end
      StringIO.new(Xcellus::_transform(input.to_json))
    end

    # Loads a XLSX file from the provided path, returning an object of type
    # Xcellus::Instance, that exposes methods for manipulating data in the
    # file
    def load(path)
      unless path.kind_of? String
        raise ArgumentError, 'Xcellus.load expects a string path'
      end
      handle = Xcellus::_load(path)
      return Xcellus::Instance.new(handle)
    end

    # with opens the provided file and passes the loaded instance to the
    # provided block. It ensures `close` is called, freeing resources.
    # Example:
    # Xcellus.with '/path/to/file' do |file|
    #   file.append ...
    #   file.replace_row ...
    #   file.save '/new/path/to/file'
    # end
    def with(path)
      raise ArgumentError, 'with requires a block' unless block_given?
      instance = load(path)
      result = yield instance
      instance.close
      return result
    end
  end

  # Instance represents a XLSX file, and exposes methods to manipulate data
  # on the file
  class Instance
    # Internal: Creates a new instance with the provided handle
    def initialize(handle)
      @handle = handle
    end

    # Searches a given sheet for a provided value in a specific column.
    # sheet_name:   Name of the sheet to lookup for `value`. Immediately returns
    #               -1 when the sheet cannot be found.
    # column_index: Index of the column to lookup for the provided value.
    # value:        Value to lookup for. Automatically converted and compared
    #               as an String.
    def find_in_column(sheet_name, column_index, value)
      unless sheet_name.kind_of? String
        raise ArgumentError, 'Invalid sheet name'
      end
      unless column_index.kind_of? Integer
        raise ArgumentError, 'Invalid column index'
      end
      Xcellus::_find_in_column(@handle, sheet_name, column_index, value.to_s)
    end

    # Replaces the row at `index` in the provided `sheet_name`.
    # sheet_name:  Name of the sheet in which the row must be replaced. Throws
    #              a StandardException when a sheet with the provided name
    #              cannot be found.
    # index:      Index of the row to be replaced.
    # value:      An array with values to be replaced. Passing `nil` prevents
    #             values of the cell in the same index from being changed.
    def replace_row(sheet_name, index, value)
      unless sheet_name.kind_of? String
        raise ArgumentError, 'Invalid sheet name'
      end
      unless index.kind_of? Integer
        raise ArgumentError, 'Invalid column index'
      end
      unless value.kind_of? Array
        raise ArgumentError, 'Invalid value: should be an array'
      end
      Xcellus::_replace_row(@handle, sheet_name, value.to_json, index)
    end

    # Close ensures that all resources allocated when `load` is called are
    # freed. This method MUST be called after you're done handling data.
    def close
      Xcellus::_close(@handle)
    end

    # Saves the current modifications to the provided path.
    def save(path)
      unless path.kind_of? String
        raise ArgumentError, 'save expects a string path'
      end

      Xcellus::_save(@handle, path)
    end

    # Appends sheets and rows to the loaded file. This method expects the same
    # structure of Xcellus::transform, with the difference that it creates (when
    # necessary) sheets, and appends data to them.
    def append(data)
      unless input.kind_of? Array
        raise ArgumentError, 'Xcellus.append only accepts Arrays'
      end
      StringIO.new(Xcellus::_transform(input.to_json))
    end
  end
end
