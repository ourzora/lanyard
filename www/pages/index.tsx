import CreateRoot from 'components/CreateRoot'
import FAQ from 'components/FAQ'

export default function CreatePage() {
  return (
    <div className="flex flex-col">
      <div className="font-bold text-2xl text-center my-10">
        Create an allow list in seconds that works across web3
      </div>
      <CreateRoot />

      <div className="h-px w-full bg-neutral-300 my-24" />

      <FAQ />

      <div className="h-px w-full bg-neutral-300 my-24" />

      <div className="flex flex-col md:flex-row gap-x-20 gap-y-10 items-center justify-center mb-16">
        <img src="/zora.png" alt="Zora logo" className="w-52" />
        <img src="/mintfun.png" alt="mint.fun logo" className="w-52" />
        <img src="/context.png" alt="Context logo" className="w-52" />
      </div>
    </div>
  )
}
