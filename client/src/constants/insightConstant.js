const VEHICLE_TYPE = ["motorcycle", "car", "bus", "truck"];

export function getInsightConstant(analyticID) {
  const analyticCode = analyticID.split("-");
  switch (analyticCode[1]) {
    case "VC": {
      const VEHICLE_COUNTING = [
        {
          label: "Total vehicle recognized today",
          realtime: true,
          keyword: "total_today"
        },
        {
          label: "Total vehicle detected last 7 days",
          realtime: false,
          keyword: "total_week"
        }
      ];

      for (let i = 0; i < VEHICLE_TYPE.length; i += 1) {
        VEHICLE_COUNTING.push({
          label: `Total ${VEHICLE_TYPE[i]} detected today`,
          realtime: true,
          keyword: VEHICLE_TYPE[i]
        });
      }
      return VEHICLE_COUNTING;
    }
    default:
      return [];
  }
}
