<!DOCTYPE html>
<html>
  <head>
    <title>beego welcome template</title>
  </head>
  <body>
{{template "block"}}
{{.hello}}
{{template "header"}}
{{template "blocks/block.tpl"}}
  </body>
</html>
