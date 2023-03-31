import CreateRoot from 'components/CreateRoot'
import About from 'components/About'
import PageTitle from 'components/PageTitle'

export default function CreatePage() {
  return (
    <div className="flex flex-col">
      <PageTitle>
        Create an allowlist in seconds that works across web3
      </PageTitle>

      <CreateRoot />

      <About />
    </div>
  )
}
