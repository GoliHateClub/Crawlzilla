package menus

var MainMenu = [][]MenuItem{
	{
		{Path: "/add_admin", IsAdmin: true, Name: "اضافه کردن ادمین"},
		{Path: "/remove_admin", IsAdmin: true, Name: "حذف کردن ادمین"},
	},
	{
		{Path: "/get_admin", IsAdmin: true, Name: "نمایش اطلاعات ادمین"},
		{Path: "/get_all_users", IsAdmin: true, Name: "نمایش همه کاربران"},
	},
	{
		{Path: "/filters", IsAdmin: false, Name: "نمایش همه فیلتر ها"},
		{Path: "/add_filter", IsAdmin: false, Name: "اضافه کردن فیلتر"},
	},
	{
		{Path: "/search", IsAdmin: false, Name: "نمایش آگهی ها"},
	},
	{
		{Path: "/add_ad", IsAdmin: true, Name: "اضافه کردن آگهی"},
		{Path: "/remove_ad", IsAdmin: true, Name: "حذف کردن آگهی"},
		{Path: "/update_ad", IsAdmin: true, Name: "ویرایش کردن آگهی"},
	},
}
