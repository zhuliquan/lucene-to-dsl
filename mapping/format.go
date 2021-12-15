package mapping

// built-in format
var BuiltInFormat = map[string][]string{
	// A formatter for the number of milliseconds since the epoch.
	// Note, that this timestamp is subject to the limits of a Java Long.MIN_VALUE and Long.MAX_VALUE.
	"epoch_millis": {""},
	// A formatter for the number of seconds since the epoch.
	// Note, that this timestamp is subject to the limits of a Java Long.MIN_VALUE and Long.
	// MAX_VALUE divided by 1000 (the number of milliseconds in a second).
	// date_optional_time or strict_date_optional_time
	// A generic ISO datetime parser, where the date must include the year at a minimum,
	// and the time (separated by T), is optional. Examples: yyyy-MM-dd'T'HH:mm:ss.SSSZ or yyyy-MM-dd.
	"epoch_second": {""},

	// A generic ISO datetime parser, where the date must include the year at a minimum, and the time (separated by T), is optional. The fraction of a second part has a nanosecond resolution. Examples: yyyy-MM-dd'T'HH:mm:ss.SSSSSSZ or yyyy-MM-dd.
	"strict_date_optional_time_nanos": {"yyyy-MM-dd'T'HH:mm:ss.SSSSSSZ", "yyyy-MM-dd"},

	// A basic formatter for a full date as four digit year, two digit month of year, and two digit day of month: yyyyMMdd.
	"basic_date": {"yyyyMMdd"},

	// A basic formatter that combines a basic date and time, separated by a T: .
	"basic_date_time": {"yyyyMMdd'T'HHmmss.SSSZ"},

	// A basic formatter that combines a basic date and time without millis, separated by a T: yyyyMMdd'T'HHmmssZ.
	"basic_date_time_no_millis": {"yyyyMMdd'T'HHmmssZ"},

	// A formatter for a full ordinal date, using a four digit year and three digit dayOfYear: yyyyDDD.
	"basic_ordinal_date": {"yyyyDDD"},

	// A formatter for a full ordinal date and time, using a four digit year and three digit dayOfYear: yyyyDDD'T'HHmmss.SSSZ.
	"basic_ordinal_date_time": {"yyyyDDD'T'HHmmss.SSSZ"},

	// A formatter for a full ordinal date and time without millis, using a four digit year and three digit dayOfYear: yyyyDDD'T'HHmmssZ.
	"basic_ordinal_date_time_no_millis": {"yyyyDDD'T'HHmmssZ"},

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit millis, and time zone offset: HHmmss.SSSZ.
	"basic_time": {"HHmmss.SSSZ"},

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset: HHmmssZ.
	"basic_time_no_millis": {"HHmmssZ"},

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit millis, and time zone off set prefixed by T: 'T'HHmmss.SSSZ.
	"basic_t_time": {"'T'HHmmss.SSSZ"},

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset prefixed by T: 'T'HHmmssZ.
	"basic_t_time_no_millis": {"'T'HHmmssZ"},

	// A basic formatter for a full date as four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx'W'wwe.
	"basic_week_date":        {"xxxx'W'wwe"},
	"strict_basic_week_date": {"xxxx'W'wwe"},

	// A basic formatter that combines a basic weekyear date and time, separated by a T: xxxx'W'wwe'T'HHmmss.SSSZ.
	"basic_week_date_time":        {"xxxx'W'wwe'T'HHmmss.SSSZ"},
	"strict_basic_week_date_time": {"xxxx'W'wwe'T'HHmmss.SSSZ"},

	// A basic formatter that combines a basic weekyear date and time without millis, separated by a T: xxxx'W'wwe'T'HHmmssZ.
	"basic_week_date_time_no_millis":        {"xxxx'W'wwe'T'HHmmssZ"},
	"strict_basic_week_date_time_no_millis": {"xxxx'W'wwe'T'HHmmssZ"},

	// A formatter for a full date as four digit year, two digit month of year, and two digit day of month: yyyy-MM-dd.
	"date":        {"yyyy-MM-dd"},
	"strict_date": {"yyyy-MM-dd"},

	// A formatter that combines a full date and two digit hour of day: yyyy-MM-dd'T'HH.
	"date_hour":        {"yyyy-MM-dd'T'HH"},
	"strict_date_hour": {"yyyy-MM-dd'T'HH"},

	// A formatter that combines a full date, two digit hour of day, and two digit minute of hour: yyyy-MM-dd'T'HH:mm.
	"date_hour_minute":        {"yyyy-MM-dd'T'HH:mm"},
	"strict_date_hour_minute": {"yyyy-MM-dd'T'HH:mm"},

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, and two digit second of minute: yyyy-MM-dd'T'HH:mm:ss.
	"date_hour_minute_second":        {"yyyy-MM-dd'T'HH:mm:ss"},
	"strict_date_hour_minute_second": {"yyyy-MM-dd'T'HH:mm:ss"},

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: yyyy-MM-dd'T'HH:mm:ss.SSS.
	"date_hour_minute_second_fraction":        {"yyyy-MM-dd'T'HH:mm:ss.SSS"},
	"strict_date_hour_minute_second_fraction": {"yyyy-MM-dd'T'HH:mm:ss.SSS"},

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: yyyy-MM-dd'T'HH:mm:ss.SSS.
	"date_hour_minute_second_millis":        {"yyyy-MM-dd'T'HH:mm:ss.SSS"},
	"strict_date_hour_minute_second_millis": {"yyyy-MM-dd'T'HH:mm:ss.SSS"},

	// A formatter that combines a full date and time, separated by a T: yyyy-MM-dd'T'HH:mm:ss.SSSZZ.
	"date_time":        {"yyyy-MM-dd'T'HH:mm:ss.SSSZZ"},
	"strict_date_time": {"yyyy-MM-dd'T'HH:mm:ss.SSSZZ"},

	// A formatter that combines a full date and time without millis, separated by a T: yyyy-MM-dd'T'HH:mm:ssZZ.
	"date_time_no_millis":        {"yyyy-MM-dd'T'HH:mm:ssZZ"},
	"strict_date_time_no_millis": {"yyyy-MM-dd'T'HH:mm:ssZZ"},

	// A formatter for a two digit hour of day: HH
	"hour":        {"HH"},
	"strict_hour": {"HH"},

	// A formatter for a two digit hour of day and two digit minute of hour: HH:mm.
	"hour_minute":        {"HH:mm"},
	"strict_hour_minute": {"HH:mm"},

	// A formatter for a two digit hour of day, two digit minute of hour, and two digit second of minute: HH:mm:ss.
	"hour_minute_second":        {"HH:mm:ss"},
	"strict_hour_minute_second": {"HH:mm:ss"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: HH:mm:ss.SSS.
	"hour_minute_second_fraction":        {"HH:mm:ss.SSS"},
	"strict_hour_minute_second_fraction": {"HH:mm:ss.SSS"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: HH:mm:ss.SSS.
	"hour_minute_second_millis":        {"HH:mm:ss.SSS"},
	"strict_hour_minute_second_millis": {"HH:mm:ss.SSS"},

	// A formatter for a full ordinal date, using a four digit year and three digit dayOfYear: yyyy-DDD.
	"ordinal_date":        {"yyyy-DDD"},
	"strict_ordinal_date": {"yyyy-DDD"},

	// A formatter for a full ordinal date and time, using a four digit year and three digit dayOfYear: yyyy-DDD'T'HH:mm:ss.SSSZZ.
	"ordinal_date_time":        {"yyyy-DDD'T'HH:mm:ss.SSSZZ"},
	"strict_ordinal_date_time": {"yyyy-DDD'T'HH:mm:ss.SSSZZ"},

	// A formatter for a full ordinal date and time without millis, using a four digit year and three digit dayOfYear: yyyy-DDD'T'HH:mm:ssZZ.
	"ordinal_date_time_no_millis":        {"yyyy-DDD'T'HH:mm:ssZZ"},
	"strict_ordinal_date_time_no_millis": {"yyyy-DDD'T'HH:mm:ssZZ"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit fraction of second, and time zone offset: HH:mm:ss.SSSZZ.
	"time":        {"HH:mm:ss.SSSZZ"},
	"strict_time": {"HH:mm:ss.SSSZZ"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset: HH:mm:ssZZ.
	"time_no_millis":        {"HH:mm:ssZZ"},
	"strict_time_no_millis": {"HH:mm:ssZZ"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit fraction of second, and time zone offset prefixed by T: 'T'HH:mm:ss.SSSZZ.
	"t_time":        {"'T'HH:mm:ss.SSSzz"},
	"strict_t_time": {"'T'HH:mm:ss.SSSzz"},

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset prefixed by T: 'T'HH:mm:ssZZ.
	"t_time_no_millis":        {"'T'HH:mm:ssZZ"},
	"strict_t_time_no_millis": {"'T'HH:mm:ssZZ"},

	// A formatter for a full date as four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx-'W'ww-e.
	"week_date":        {"xxxx-'W'ww-e"},
	"strict_week_date": {"xxxx-'W'ww-e"},

	// A formatter that combines a full weekyear date and time, separated by a T: xxxx-'W'ww-e'T'HH:mm:ss.SSSZZ.
	"week_date_time":        {"xxxx-'W'ww-e'T'HH:mm:ss.SSSZZ"},
	"strict_week_date_time": {"xxxx-'W'ww-e'T'HH:mm:ss.SSSZZ"},

	// A formatter that combines a full weekyear date and time without millis, separated by a T: xxxx-'W'ww-e'T'HH:mm:ssZZ.
	"week_date_time_no_millis":        {"xxxx-'W'ww-e'T'HH:mm:ssZZ"},
	"strict_week_date_time_no_millis": {"xxxx-'W'ww-e'T'HH:mm:ssZZ"},

	// A formatter for a four digit weekyear: xxxx.
	"weekyear":        {"xxxx"},
	"strict_weekyear": {"xxxx"},

	// A formatter for a four digit weekyear and two digit week of weekyear: xxxx-'W'ww.
	"weekyear_week":        {"xxxx-'W'ww"},
	"strict_weekyear_week": {"xxxx-'W'ww"},

	// A formatter for a four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx-'W'ww-e.
	"weekyear_week_day":        {"xxxx-'W'ww-e"},
	"strict_weekyear_week_day": {"xxxx-'W'ww-e"},

	// A formatter for a four digit year and two digit month of year: yyyy-MM.
	"year_month":        {"yyyy-MM"},
	"strict_year_month": {"yyyy-MM"},

	// A formatter for a four digit year: yyyy.
	"year":        {"yyyy"},
	"strict_year": {"yyyy"},

	// A formatter for a four digit year, two digit month of year, and two digit day of month: yyyy-MM-dd.
	"year_month_day":        {"yyyy-MM-dd"},
	"strict_year_month_day": {"yyyy-MM-dd"},
}
