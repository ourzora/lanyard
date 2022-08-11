import CreateRoot from 'components/CreateRoot'
import FAQ from 'components/FAQ'

export default function CreatePage() {
  return (
    <div className="flex flex-col mb-24">
      <div className="font-bold text-2xl text-center my-10">
        Create an allow list in seconds that works across web3
      </div>

      <CreateRoot />

      <div
        className="h-px bg-neutral-200 my-16"
        style={{ marginLeft: '-100vw', marginRight: '-100vw' }}
      />

      <div className="flex flex-col gap-y-6">
        <div className="text-xl">Built in collaboration with</div>
        <div className="flex flex-col sm:flex-row gap-x-10 gap-y-10 items-center">
          <img src="/zora.png" alt="Zora logo" className="w-24" />
          <img src="/mintfun.png" alt="mint.fun logo" className="w-24" />
          <img src="/context.png" alt="Context logo" className="w-24" />
        </div>
      </div>

      <div
        className="h-px bg-neutral-200 my-16"
        style={{ marginLeft: '-100vw', marginRight: '-100vw' }}
      />

      <FAQ />
    </div>
  )
}
