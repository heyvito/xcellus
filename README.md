# Xcellus

[![Gem Version](https://badge.fury.io/rb/xcellus.svg)](https://badge.fury.io/rb/xcellus)

**Xcellus** is an ugly XLSX generator. It leverages a native extension to create
XLSX files. It is said to be an 'ugly' generator because it lacks all styling
mechanisms to make your XLSX pretty (such as fonts, borders, and the likes).

## Installation

> **WARNING**: Xcellus is targeted for Linux and macOS (both amd64). It won't work on other architectures and support for them is not planned.

Add this line to your application's Gemfile:

```ruby
gem 'xcellus'
```

And then execute:

    $ bundle

Or install it yourself as:

    $ gem install xcellus

## Usage

Xcellus transforms a list of hashes (acting as worksheets) into a XLSX workbook:

```ruby
require 'xcellus'
io = Xcellus.transform [
  {
    title: 'Brian\'s Worksheet',
    headers: [:Artist, :Track, :Playcount],
    rows: [
        ['Metallica',               'Hero of the Day',      242],
        ['Metallica',               'The Shortest Straw',   186],
        ['Queens of the Stone Age', 'My God Is the Sun',    276],
        ['Queens of the Stone Age', 'I Sat by the Ocean',   270],
        ['Gorillaz',                'On Melancholy Hill',   203],
        ['Gorillaz',                'Kids With Guns',       184],
    ]
  }
]
# => #<StringIO:0x000000000245cc78>

File.open('BrianosMusic.xlsx', 'wb') do |f|
  f.write(io.string)
end
# => 6377
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/victorgama/xcellus. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## Code of Conduct

Everyone interacting in the Xcellus project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](https://github.com/victorgama/xcellus/blob/master/CODE_OF_CONDUCT.md).

## License

```
The MIT License (MIT)

Copyright (c) 2018 Victor Gama

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
