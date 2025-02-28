import { Hero } from "@/components/ui/home/hero";
import { About } from "@/components/ui/home/about";
import { GettingStarted } from "@/components/ui/home/getting-started";

export default function Home() {
  return (
    <main>
      <div className="flex flex-col justify-center justify-items-center gap-y-10">
        <section id="/" className="pt-20">
          <Hero />
        </section>
        <section id="about">
          <About />
        </section>
        <section>
          <GettingStarted />
        </section>
      </div>
    </main>
  );
}
