import { SVGProps } from 'react'

export default function Logo({
  height = 68,
  width = 62,
  className,
}: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 62 68"
      width={width}
      height={height}
      className={className}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M62 31C62 26.929 61.1982 22.8979 59.6403 19.1368C58.0824 15.3757 55.7989 11.9583 52.9203 9.07969C50.0417 6.20107 46.6243 3.91763 42.8632 2.35973C39.1021 0.801838 35.071 0 31 0C26.929 0 22.8979 0.801839 19.1368 2.35974C15.3757 3.91763 11.9583 6.20107 9.07969 9.07969C6.20107 11.9583 3.91763 15.3757 2.35973 19.1368C0.801838 22.8979 0 26.929 0 31C0 35.071 0 68 0 68H10C10 68 10 42.598 10 31C10 19.402 19.402 10 31 10C42.598 10 52 19.402 52 31C52 42.598 52 68 52 68H62C62 68 62 35.071 62 31ZM31 16C33.7614 16 36 18.2386 36 21V26H26V21C26 18.2386 28.2386 16 31 16ZM26 26V62V63V68H21H16V63V62V36C16 30.4772 20.4772 26 26 26ZM46 62L46 36C46 30.4772 41.5228 26 36 26V62L36 63L36 63.0062V68H41H46V63V62Z"
        className="fill-current"
      />
    </svg>
  )
}
