# For next release
  * **Raphael Pour**
    * render: 
      * disable favicon to avoid client request
      * add markdown footnote support
      * **bugfix** navigation: render next/previous links properly and ignore
        static posts

*Not released yet*

# Minor Release v0.4.0 (2022-07-30)
  * **Raphael Pour**
    * **bugfix** cmd/list: fix file count condition to check if a post directory has enough files
    * post: add static posts that are unlisted
    * cmd/list: add 'static' column

*Released by Raphael Pour <info@raphaelpour.de>*

# Minor Release v0.3.0 (2022-03-23)
  * **Raphael Pour**
    * render: add navigation with previous+next post and home
  * **Tch1b0**
    * render: add syntax-highlighting for code-blocks
    

*Released by Raphael Pour <info@raphaelpour.de>*

# Minor Release v0.2.0 (2022-01-13)
  * **Raphael Pour**
    * **bugfix** metadata: fix typo in status key
    * render: 
      * separate posts with `<hr>`
      * sort post from latest to oldest
      * change from single- to multi-page site
      * fix code width
      * generate rss feed
      * add image support

*Released by Raphael Pour <info@raphaelpour.de>*

# Patch Release v0.1.3 (2020-11-10)
  * **Raphael Pour**
    * **bugfix** Fix force and output flag
    * Rendering: Add creation date to post
    * Cmd: Add list command
    * Add Status: 
      * Track draft/public posts
      * Only render public posts
      * Add cmd draft/publish to set a post's state
      * Adjust cmds to handle status properly

*Released by Raphael Pour <info@raphaelpour.de>*

# Patch Release v0.1.2 (2020-10-04)
  * **Raphael Pour**
    * Update ci

*Released by Raphael Pour <info@raphaelpour.de>*

# Patch Release v0.1.1 (2020-10-04)
  * **Raphael Pour**
    * Add travis deployment

*Released by Raphael Pour <info@raphaelpour.de>*

# Minor Release v0.1.0 (2020-10-04)
  * **Raphael Pour**
    * Add commands:
      * Init
      * Add
      * Update
      * Render

*Released by Raphael Pour <info@raphaelpour.de>*
