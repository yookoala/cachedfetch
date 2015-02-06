crawler
=======

A library to fetch and cache HTTP results.

What it does?
-------------

This is a library to fetch and cache HTTP results. It allows easy comparison between different version of the same site. Suitable for crawler

You may extend the way to diff versions with the `Diff` interface. That means you decide what it means to be different and how you want to keep revisions of Raw HTTP results.

The default cache uses sql database. You may extend by implementing the `Cacher` interface.

Licence
-------

This file is part of crawler.

crawler is free software: you can redistribute it and/or modify it under the terms of the GNU Lesser General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

crawler is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public Licensefor more details.

You should have received a copy of the GNU Lesser General Public License along with crawler. If not, see http://www.gnu.org/licenses/lgpl.html.
