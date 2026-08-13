export const frequency1Pattern = '^(NotKnown)$|^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlDay:((0[2-9])|([1-2][0-9])|3[0-1]))$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$';
export const frequency1Regex = new RegExp(frequency1Pattern);
export const pointInTimePattern = '^(-[0-9]|[0-9]{1,2})$';
export const pointInTimeRegex = new RegExp(pointInTimePattern);
export const maxCountPerPeriod = 2147483647;
export const obFrequency6Codes = [
  'ADHO',
  'YEAR',
  'DAIL',
  'FRTN',
  'INDA',
  'MNTH',
  'QURT',
  'MIAN',
  'WEEK',
  'WODL',
  'FOWK',
  'TWMH',
  'FOMH',
  'FIMH',
  'ALMH',
  'NONE',
  'LWMH',
  'LXMH',
  'TWYR',
];

export const hasValue = value => value !== undefined && value !== null && value !== '';
export const isFrequency1 = value => hasValue(value) && frequency1Regex.test(value);
export const isOBFrequency6Code = value => obFrequency6Codes.includes(value);
export const isPointInTime = value => hasValue(value) && pointInTimeRegex.test(value);
export const isCountPerPeriod = value => Number.isInteger(value) && value > 0 && value <= maxCountPerPeriod;

export const validatePaymentFrequency = value => !hasValue(value) || isFrequency1(value);

export const validateV4StandingOrderFrequency = (frequency) => {
  const errors = [];
  const value = frequency || {};
  const hasPointInTime = hasValue(value.PointInTime);
  const hasCountPerPeriod = hasValue(value.CountPerPeriod);
  const typeIsCode = isOBFrequency6Code(value.Type);
  const typeIsFrequency1 = isFrequency1(value.Type);

  if (!hasValue(value.Type)) {
    errors.push('v4_standing_order_frequency.Type empty');
  } else if (!typeIsCode && !typeIsFrequency1) {
    errors.push('v4_standing_order_frequency.Type invalid');
  }
  if (typeIsFrequency1 && !typeIsCode && (hasPointInTime || hasCountPerPeriod)) {
    errors.push('v4_standing_order_frequency legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod');
  }
  if (hasPointInTime && !isPointInTime(value.PointInTime)) {
    errors.push('v4_standing_order_frequency.PointInTime must be numeric text up to two characters, including negative single-digit values');
  }
  if (hasCountPerPeriod && !isCountPerPeriod(value.CountPerPeriod)) {
    errors.push('v4_standing_order_frequency.CountPerPeriod must be a positive Int32');
  }
  if (hasPointInTime && hasCountPerPeriod) {
    errors.push('v4_standing_order_frequency.PointInTime and CountPerPeriod are mutually exclusive');
  }

  return errors;
};
