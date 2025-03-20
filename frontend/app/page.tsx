import Link from "next/link";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { ThemeToggle } from "@/components/theme-toggle";
import {
  CloudIcon,
  ArrowRightIcon,
  ShieldCheckIcon,
  RefreshCwIcon,
  FolderIcon,
  ZapIcon,
} from "lucide-react";
import { FaDropbox, FaGoogleDrive } from "react-icons/fa6";
import { IconType } from "react-icons/lib";
// import { GoogleDriveIcon, DropboxIcon } from "@/components/icons";

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col">
      {/* Navigation */}
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-16 items-center justify-between">
          <div className="flex items-center gap-2">
            <CloudIcon className="h-6 w-6 text-primary" />
            <span className="text-xl font-bold">Cloudmesh</span>
          </div>

          <nav className="hidden md:flex gap-6">
            <Link
              href="#features"
              className="text-sm font-medium hover:text-primary"
            >
              Features
            </Link>
            <Link
              href="#how-it-works"
              className="text-sm font-medium hover:text-primary"
            >
              How It Works
            </Link>
            <Link
              href="#providers"
              className="text-sm font-medium hover:text-primary"
            >
              Providers
            </Link>
            <Link
              href="#faq"
              className="text-sm font-medium hover:text-primary"
            >
              FAQ
            </Link>
          </nav>

          <div className="flex items-center gap-2">
            <ThemeToggle />
            <Link href="/auth">
              <Button size="sm">Sign Up</Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section className="container py-12 md:py-24 lg:py-32 grid gap-6 md:grid-cols-2 md:gap-10">
          <div className="flex flex-col justify-center space-y-4">
            <div className="space-y-2">
              <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl">
                All Your Cloud Storage in One Place
              </h1>
              <p className="text-muted-foreground md:text-xl">
                Access, manage, and transfer files between Google Drive and
                Dropbox seamlessly
              </p>
            </div>
            <div className="flex flex-col gap-2 min-[400px]:flex-row">
              <Link href="/auth">
                <Button size="lg" className="gap-2">
                  Get Started
                  <ArrowRightIcon className="h-4 w-4" />
                </Button>
              </Link>
              <Link href="#how-it-works">
                <Button variant="outline" size="lg">
                  Learn More
                </Button>
              </Link>
            </div>
          </div>
          <div className="flex items-center justify-center">
            <div className="relative w-full max-w-[500px] aspect-[4/3]">
              <Image
                src="/placeholder.svg?height=600&width=800"
                alt="Cloudmesh interface"
                width={600}
                height={450}
                className="rounded-lg shadow-lg object-cover"
              />
              <div className="absolute -bottom-4 -right-4 bg-background rounded-lg shadow-lg p-3 border">
                <div className="flex items-center gap-2">
                  <FaGoogleDrive className="h-6 w-6" />
                  <ArrowRightIcon className="h-4 w-4 text-muted-foreground" />
                  <FaDropbox className="h-6 w-6" />
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section id="features" className="bg-muted/50 py-12 md:py-24 lg:py-32">
          <div className="container space-y-12">
            <div className="text-center space-y-4">
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
                Key Features
              </h2>
              <p className="text-muted-foreground md:text-xl max-w-[700px] mx-auto">
                Cloudmesh simplifies how you work with files across multiple
                cloud storage providers
              </p>
            </div>

            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              <FeatureCard
                Icon={FolderIcon}
                title="Unified Access"
                description="Browse and manage all your cloud files in a single, intuitive interface"
              />
              <FeatureCard
                Icon={ArrowRightIcon}
                title="Easy Transfers"
                description="Move or copy files between cloud providers with just a few clicks"
              />
              <FeatureCard
                Icon={RefreshCwIcon}
                title="No Local Storage Required"
                description="Work with your files without downloading them to your device"
              />
              <FeatureCard
                Icon={ShieldCheckIcon}
                title="Secure Authentication"
                description="Industry-standard OAuth ensures your credentials are never exposed"
              />
            </div>
          </div>
        </section>

        {/* How It Works Section */}
        <section id="how-it-works" className="py-12 md:py-24 lg:py-32">
          <div className="container space-y-12">
            <div className="text-center space-y-4">
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
                How It Works
              </h2>
              <p className="text-muted-foreground md:text-xl max-w-[700px] mx-auto">
                Get started with Cloudmesh in three simple steps
              </p>
            </div>

            <div className="grid gap-10 md:grid-cols-3">
              <StepCard
                number="1"
                title="Connect Your Accounts"
                description="Securely link your Google Drive and Dropbox accounts with OAuth"
                image="/placeholder.svg?height=200&width=300"
              />
              <StepCard
                number="2"
                title="Browse All Your Files"
                description="View and manage files from all your cloud providers in one interface"
                image="/placeholder.svg?height=200&width=300"
              />
              <StepCard
                number="3"
                title="Transfer With Ease"
                description="Move or copy files between providers without downloading them first"
                image="/placeholder.svg?height=200&width=300"
              />
            </div>
          </div>
        </section>

        {/* Providers Section */}
        <section id="providers" className="bg-muted/50 py-12 md:py-24 lg:py-32">
          <div className="container space-y-12">
            <div className="text-center space-y-4">
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
                Supported Providers
              </h2>
              <p className="text-muted-foreground md:text-xl max-w-[700px] mx-auto">
                Connect your favorite cloud storage services, with more coming
                soon
              </p>
            </div>

            <div className="flex flex-wrap justify-center gap-8 md:gap-16">
              <div className="flex flex-col items-center gap-4">
                <FaGoogleDrive className="h-20 w-20" />
                <span className="text-lg font-medium">Google Drive</span>
              </div>
              <div className="flex flex-col items-center gap-4">
                <FaDropbox className="h-20 w-20" />
                <span className="text-lg font-medium">Dropbox</span>
              </div>
            </div>

            <div className="text-center">
              <p className="text-muted-foreground italic">
                More providers coming soon
              </p>
            </div>
          </div>
        </section>

        {/* FAQ Section */}
        <section id="faq" className="py-12 md:py-24 lg:py-32">
          <div className="container space-y-12">
            <div className="text-center space-y-4">
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
                Frequently Asked Questions
              </h2>
              <p className="text-muted-foreground md:text-xl max-w-[700px] mx-auto">
                Find answers to common questions about Cloudmesh
              </p>
            </div>

            <div className="mx-auto max-w-3xl">
              <Accordion type="single" collapsible className="w-full">
                <AccordionItem value="item-1">
                  <AccordionTrigger>Is my data secure?</AccordionTrigger>
                  <AccordionContent>
                    {
                      "Yes, your data remains secure. Cloudmesh doesn't store your\
                    files on our servers. We use secure API connections to your\
                    cloud providers and industry-standard OAuth for\
                    authentication. Your credentials are never exposed to our\
                    system."
                    }
                  </AccordionContent>
                </AccordionItem>
                <AccordionItem value="item-2">
                  <AccordionTrigger>Do you store my files?</AccordionTrigger>
                  <AccordionContent>
                    {
                      "No, Cloudmesh doesn't store your files. We provide a unified\
                    interface to access and manage your files directly from your\
                    cloud providers. Files remain stored on your existing cloud\
                    storage accounts."
                    }
                  </AccordionContent>
                </AccordionItem>
                <AccordionItem value="item-3">
                  <AccordionTrigger>
                    How does authentication work?
                  </AccordionTrigger>
                  <AccordionContent>
                    Cloudmesh uses OAuth 2.0, the industry standard for secure
                    authentication. This means you log in directly with your
                    cloud provider (Google, Dropbox) and grant Cloudmesh limited
                    access permissions. We never see or store your passwords.
                  </AccordionContent>
                </AccordionItem>
                <AccordionItem value="item-4">
                  <AccordionTrigger>
                    What file operations are supported?
                  </AccordionTrigger>
                  <AccordionContent>
                    Cloudmesh supports viewing, downloading, uploading, copying,
                    moving, and deleting files. You can also transfer files
                    directly between cloud providers without downloading them to
                    your device first, saving time and bandwidth.
                  </AccordionContent>
                </AccordionItem>
                <AccordionItem value="item-5">
                  <AccordionTrigger>
                    Can I use Cloudmesh on mobile devices?
                  </AccordionTrigger>
                  <AccordionContent>
                    {
                      "Yes, Cloudmesh is fully responsive and works on smartphones\
                    and tablets. Access your cloud files from any device with a\
                    web browser. We're also working on dedicated mobile apps for\
                    iOS and Android."
                    }
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="bg-muted/50 py-12 md:py-24 lg:py-32">
          <div className="container flex flex-col items-center text-center space-y-4">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
              Ready to simplify your cloud storage?
            </h2>
            <p className="text-muted-foreground md:text-xl max-w-[700px]">
              Join thousands of users who have streamlined their cloud file
              management with Cloudmesh
            </p>
            <Link href="/auth" className="mt-4">
              <Button size="lg" className="gap-2">
                Get Started
                <ZapIcon className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t py-6 md:py-8">
        <div className="container flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-2">
            <CloudIcon className="h-5 w-5 text-primary" />
            <span className="text-lg font-semibold">Cloudmesh</span>
            <span className="text-sm text-muted-foreground ml-2">
              All your cloud storage in one place
            </span>
          </div>

          <div className="flex gap-4 text-sm text-muted-foreground">
            <Link href="#" className="hover:underline">
              Terms of Service
            </Link>
            <Link href="#" className="hover:underline">
              Privacy Policy
            </Link>
          </div>

          <div className="text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} Cloudmesh. All rights reserved.
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({
  Icon,
  title,
  description,
}: {
  Icon: IconType;
  title: string;
  description: string;
}) {
  return (
    <div className="group relative overflow-hidden rounded-lg border bg-background p-6 hover:shadow-md transition-all duration-300 hover:-translate-y-1">
      <div className="space-y-2">
        <div className="text-primary">
          <Icon className="h-20 w-20" />
        </div>
        <h3 className="font-bold">{title}</h3>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}

function StepCard({
  number,
  title,
  description,
  image,
}: {
  number: string;
  title: string;
  description: string;
  image: string;
}) {
  return (
    <div className="flex flex-col items-center text-center space-y-4">
      <div className="relative">
        <div className="absolute -left-3 -top-3 flex h-10 w-10 items-center justify-center rounded-full bg-primary text-lg font-bold text-primary-foreground">
          {number}
        </div>
        <Image
          src={image || "/placeholder.svg"}
          alt={title}
          width={300}
          height={200}
          className="rounded-lg border object-cover aspect-[3/2]"
        />
      </div>
      <h3 className="text-xl font-bold">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}
