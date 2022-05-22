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

export async function getAllData ( minLat:number, maxLat:number, minLng:number, maxLng:number, zoom:number , currentType:number, startDate: Date, endDate:Date) {
    const startTs = startDate.getTime() / 1000
    const endTs = endDate.getTime()/1000
    const res = await axios.get(
        process.env.NEXT_PUBLIC_SERVER_HOST+"/cubes/time", 
        {params: 
            {minLat: minLat, 
                maxLat:maxLat, 
                minLng:minLng,
                maxLng:maxLng, 
                zoom:zoom,
                type: currentType,
                startTime: startTs,
                endTime: endTs
            }});
    return res;
}

export async function getTypes (){
    const res = await axios.get(process.env.NEXT_PUBLIC_SERVER_HOST+"/types");
    return res;
}