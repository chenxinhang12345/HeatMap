export function convertNanoCubeBoundsToNSEW(nanoCubeBounds: any ){
    const {bounds, count} = nanoCubeBounds;
    const {lng, lat, width, height} = bounds;
    return {
        north: lat,
        south: lat - height ,
        east: lng+width,
        west: lng
    }
}