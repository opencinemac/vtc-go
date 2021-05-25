/*
Package rate offers methods and types for parsing and constructing framerate values.

Parsing and NTSC Framrates

When parsing NTSC Framerates, the playback speed will be rounded to the nearest valid
NTSC value of n/1001, so '24' and '23.98' will both get coerced to 24000/1001.

Non-NTSC Framerates will be left as-is, since there is no distinction between the
framerate and the timebase. Attempting to parse a float value will result in an error
for non-ntsc Framerates, as floats are not precise enough to know with certainty that
the correct, TRUE desired rate is being used.
*/
package rate
