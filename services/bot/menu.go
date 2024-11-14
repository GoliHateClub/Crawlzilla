package bot

type MenuItem struct {
	Path    string
	Name    string
	IsAdmin bool
}

var MainMenu = [][]MenuItem{
	{
		{Path: "/add_admin", IsAdmin: true, Name: "اضافه کردن ادمین"},
		{Path: "/remove_admin", IsAdmin: true, Name: "حذف کردن ادمین"},
	},
	{
		{Path: "/add_ad", IsAdmin: true, Name: "اضافه کردن آگهی"},
		{Path: "/remove_ad", IsAdmin: true, Name: "حذف کردن آگهی"},
		{Path: "/update_ad", IsAdmin: true, Name: "ویرایش کردن آگهی"},
	},
	{
		{Path: "/see_all_ads", IsAdmin: false, Name: "دیدن همه آگهی ها"},
	},
}
