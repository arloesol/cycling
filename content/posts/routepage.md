---
title: "Route Page"
subtitle: ""
date: 2022-11-18T13:14:04+01:00
draft: false
author: ""
authorLink: ""
authorEmail: ""
description: ""
keywords: "about"
comment: false
weight: 999

tags:
- help
categories:
- help

hiddenFromHomePage: false
hiddenFromSearch: false

summary: ""
resources:
- name: featured-image
  src: featured-image.jpg
- name: featured-image-preview
  src: featured-image-preview.jpg

toc:
  enable: true
math:
  enable: false
lightgallery: false
---

## Header info

The start of the route page contains a title and subtitle followed by some page info
<br><br>

example:

included in {{% icon class="fa-regular fa-calendar-alt fa-fw" %}} category1 {{% icon class="fa-regular fa-calendar-alt fa-fw" %}} category2

{{% icon class="fa-regular fa-calendar-alt fa-fw" %}} 2022-11-18 {{% icon class="fa-solid fa-pencil-alt fa-fw" %}} 82 words {{% icon class="fa-regular fa-clock fa-fw" %}} One minute {{% icon class="fa-solid fa-circle-info" %}} src {{% icon class="fa-solid fa-book" %}} ext

Main: {{% icon class="fa-solid fa-route" %}} 34 km {{% icon class="fa-solid fa-download" %}} gpx {{% icon class="fa-solid fa-bicycle" %}} 1.4 {{% icon class="fa-solid fa-arrow-up" %}} 220m {{% icon class="fa-solid fa-arrow-down" %}} 220m {{% icon class="fa-solid fa-angle-double-up" %}}  7.2% {{% icon class="fa-solid fa-angle-double-down" %}} -5.2%
<br><br>

- "included in" list of categories the route page is part of
- date of route page, number of words, reading duration, link to website info page (src), link to external route page (ext)
- for each gpx track related to this route
  - length in km, link to gpx file, effort level estimate, meters up, meters down, max slope %, min slope %
  - effort level is 1 for each 30km flat cycling route + some extra when it's hilly

## Route Images

Some images related to the routes may be shown in this section

## Route Content

Some information about the route is generally shown after the Route Images

The content and chapters depend on the route and the source of this route

## GPX Route

Information about the Gpx track(s) is shown here and also a map of the gpx tracks

In case there are several gpx tracks you can click on the track info to select which track is shown on the map
<br><br>

example:

Main: {{% icon class="fa-solid fa-route" %}} 34 km {{% icon class="fa-solid fa-download" %}} gpx {{% icon class="fa-solid fa-bicycle" %}} 1.4 {{% icon class="fa-solid fa-arrow-up" %}} 220m {{% icon class="fa-solid fa-arrow-down" %}} 220m {{% icon class="fa-solid fa-angle-double-up" %}}  7.2% {{% icon class="fa-solid fa-angle-double-down" %}} -5.2% {{% icon class="fa-solid fa-mountain" %}} 150m {{% icon class="fa-solid fa-water" %}} 55m

slopes: ±2%: 2,203m ±4%: 513m ±6%: 60m ±8%: 0m >9%: 0m 
<br><br>

- for each gpx track related to this route
  - length in km, link to gpx file, effort level estimate, meters up, meters down, max slope %, min slope %, max height in m and min height in m
  - effort level is 1 for each 30km flat cycling route + some extra when it's hilly
  - length for a set of slope %s
    - 2%: length in m of route with slope between 1% and 3%
    - 4%: length in m of route with slope between 1% and 3%
    - ...
    - \>9%: length in m of route with slope >9%

## Footer

The footer of the page contains

- a link to github to edit this page
- links to social media sites:app to share the page link
- list of tags linked to the route page
- links to some other linked pages
- comment system to leave comments on the page (you'll need a free github account for this)
