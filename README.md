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

It also can modify files:

```ruby
require 'xcellus'
Xcellus.with 'BrianosMusic.xlsx' do |xls|
  # Let's find and replace 'I Sat by the Ocean' with something else:
  row_index = xls.find_in_column 'Brian\'s Worksheet', 1, 'I Sat by the Ocean'
  # row_index == 4
  xls.replace_row 'Brian\'s Worksheet', row_index, [nil, 'The Evil has Landed']
  # Now, the row ['Queens of the Stone Age', 'I Sat by the Ocean',   270],
  # will be composed by ['Queens of the Stone Age', 'The Evil has Landed',   270],
  # Let's also add a new row
  xls.append [
    { title: 'Brian\'s Worksheet', rows: [['Queen', 'Spread Your Wings', 27]] }
  ]
  # Now, Brian's Worksheet has the following data:
  # | Artist                  | Track               | Playcount |
  # | Metallica               | Hero of the Day     | 242       |
  # | Metallica               | The Shortest Straw  | 186       |
  # | Queens of the Stone Age | My God Is the Sun   | 276       |
  # | Queens of the Stone Age | The Evil has Landed | 270       |
  # | Gorillaz                | On Melancholy Hill  | 203       |
  # | Gorillaz                | Kids With Guns      | 184       |
  # | Queen                   | Spread Your Wings   | 27        |
  # And we can save a copy:
  xls.save('Copy of BrianosMusic.xlsx')
end
# Here, all resources are freed.
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
