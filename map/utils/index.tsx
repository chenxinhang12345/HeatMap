import axios from 'axios';
export function convertNanoCubeBoundsToNSEW(nanoCubeBounds: any ){
    const {n,s,e,w,opacity} = nanoCubeBounds;
    return {
        rectBounds:{
        north: n,
        south: s,
        east: e,
        west: w
        },
        rectOpacity: opacity
    }
}

export async function getAllData ( minLat:number, maxLat:number, minLng:number, maxLng:number, zoom:number , currentType:number) {
    const res = await axios.get(
        process.env.NEXT_PUBLIC_SERVER_HOST+"/cubes", 
        {params: 
            {minLat: minLat, 
                maxLat:maxLat, 
                minLng:minLng,
                maxLng:maxLng, 
                zoom:zoom,
                type: currentType
            }});
    return res;
}

export async function getTypes (){
    const res = await axios.get(process.env.NEXT_PUBLIC_SERVER_HOST+"/types");
    return res;
}