# How to Use ROI Lib

1. To be able to draw inside the canvas, make sure to copy the give `canvas` as an id to the component.
2. Make sure to copy all the `#canvas` styling, image for preview can be but inside of the canvas element as children component. for example:
```
<div id="canvas>
	<img src={source} />
</div>
```
3. call `initialCanvas` function on windows onload, or when there are any ROI Type change

## Function Explanation
### initialCanvas ( callback, roi_type )
`initialCanvas` needs two parameters, namely for `callback` and `roi_type` respectively. `callback` will be called after region or line created, returning object with values:
```
{
	points:  array_of_coordinates,
	area:  name_of_area,
	color:  color_of_lines,
	lineNumber: index_of_roi (starts from 1)
}
```
### onResetROI ( callback )
this function will remove all the lines. accepting `callback` params for user to be able to call their own function when resetting all lines and region

### onDeleteLine ( id ), onDeleteRegion (id)
receive `id` as params. this `id` can be filled with `lineNumber` from the `initialCanvas` callback.
`onDeleteLine` will be used to delete lines, while `onDeleteRegion` will be used to delete area/region.

### onReverseLine ( id )
receive `id` as params. this `id` can be filled with `lineNumber` from the `initialCanvas` callback.
only applicable to `roi_type` **line**. will reverse the line direction
