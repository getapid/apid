+++
title = "CLI"
description = "How to get going with the APId CLI"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h2(text="How to install") }}

Head over to the [Github releases page](https://github.com/getapid/apid-cli/releases/) and choose your platform and architecture.
<br><br>
This should get you an archive that contains the binary. For unix follow below are the instructions:
<br><br>
```sh
tar -xzf apid-*.tar.gz
chmod u+x apid
./apid version
```