Given the src dir `-src xxx`, for all images, if width > height, cut into two equal parts, or keep it original. Copy all the results to the dst dir `-dst`.

```
cut-comic -src a -dst b -rtl

src: a, dst: b, rtl=true
1280x1877 - a/25 (1).jpg
---- Skip
480x640 - a/4af5858d30fb9688ed439c0d44110f4e_480.jpg
---- Skip
901x621 - a/s1.png
---- b/s1_1.png
---- b/s1_0.png
1280x720 - a/sample.png
---- b/sample_1.png
---- b/sample_0.png
1280x1877 - a/25.jpg
---- Skip
1280x720 - a/a.png
---- b/a_1.png
---- b/a_0.png
```