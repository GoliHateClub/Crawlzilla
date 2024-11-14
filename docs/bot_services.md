## Bot services

## Menu
<!-- TOC -->
  * [Bot services](#bot-services)
  * [Menu](#menu)
    * [User](#user)
    * [Filter](#filter)
    * [Search](#search)
    * [Bookmarks](#bookmarks)
    * [Ads](#ads)
    * [Watchlist](#watchlist)
<!-- TOC -->

### User
- Super User
  - `CreateAdmin(id int64)` `bool`, `error`
    > Get admin user id and create a super admin in database
  - `GetAdminInfo(id int64)` `User`, `error`
    > Get admin id and return its user info
  - `RemoveAdmin(id int64)` `User`, `error`
    > Get admin id and remove it from admins
  - `GetAllUsersInfo()` `struct{data []User, pages int, page int}`, `error`
    > Get all users info in system with pagination
    
### Filter
- Super User, Admins 
  - `GetAllFilters(role ADMIN|SUPER_ADMIN|NORMAL_USER)` `struct{data []FilterData, pages int, page int}`, `error`
    > Get all filer data. Filter Data is list of filters. Admins only can see filters info but super admins can see which user created that filter too
- All Users
  - `CreateOrUpdateFilter(filter Filter)` `bool`, `error`
    > Create new filter.
  - `RemoveFilter(id string)` `bool`, `error`
    > Remove filter by uuid

### Search
- All Users
  - `SearchInAds(id string)` `struct{data []{title string, image string, id string, is_bookmard book}, pages int, page int}`, `error`
    > Paginated search in all ads by filter id. `filter_count++`
    
### Bookmarks
- All Users
  - `AddBookmark(ad_id string, user_id string)` `bool`, `error`
    > Get Ad and User id. Add Ad to bookmark list.
  - `DeleteBookmark(id string)` `bool`, `error`
    > Delete bookmark by id.
  - `GetAllBookmark()` `[]{title string, image string, id string}`, `error`
    > Get all bookmarks info
    
### Ads
- Super Users
  - `CreateAd(ad Ad)` `bool`, `error`
    > Create new ad.
  - `UpdateAd(ad Ad)` `bool`, `error`
    > Update an ad.
  - `RemoveAd(id string)` `bool`, `error`
    > Remove ad by id.
  - `GetAdInfo(id stirng)` `Ad`, `error`
    > Get ad info by id. `visit_count++`

### Watchlist
> //TODO