import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import relativeTime from "dayjs/plugin/relativeTime";
import utc from "dayjs/plugin/utc";

dayjs
  .extend(timezone)
  .extend(utc)
  .extend(relativeTime);

export const dateLocalFormat = "dd-MM-yyyy HH:mm";
export const dateInterFormat = "YYYY-MM-DD HH:mm:ss";
export const eventFormat = "dddd DD/MM/YYYY - HH:mm";
export const notifFormat = "HH:mm";
export const notifDateFormat = "DD/MM";
export const enrollmentDateFormat = "DD MMM YYYY - hh:mm";

const localOffset = Intl.DateTimeFormat().resolvedOptions().timeZone;

export function dateToIso(time) {
  const dateIso = dayjs(time)
    .tz(localOffset)
    .format("YYYY-MM-DDTHH:mm:ssZ");
  return dateIso;
}

export function isoToDate(time) {
  const isoDate = new Date(dayjs(time).format(dateInterFormat));
  return isoDate;
}

export function getPastTime(day) {
  const pastDate = new Date(
    new Date().setDate(new Date().getDate() - day)
  ).setHours(23, 59, 59, 59);
  return pastDate;
}

export function eventDate(time) {
  return dayjs(time).format(eventFormat);
}

export function notifTime(time) {
  return dayjs(time).format(notifFormat);
}

export function enrollmentTime(time) {
  const enrollTime = dayjs(time)
    .tz(localOffset)
    .format(enrollmentDateFormat);
  return enrollTime;
}

export function getTimeZone() {
  return dayjs.tz.guess(true);
}

export function notifDate(time) {
  return dayjs(time).format(notifDateFormat);
}

export function isToday(time) {
  const today = dateToIso(new Date());
  const formattedToday = dayjs(today).format(notifDateFormat);
  const formattedTime = dayjs(time).format(notifDateFormat);
  return formattedTime === formattedToday;
}
