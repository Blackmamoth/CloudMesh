"use client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";
import { MainLayout } from "@/components/layouts/main-layout";

export default function SettingsPage() {
  return (
    <MainLayout>
      <div className="flex flex-col gap-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight mb-2">Settings</h1>
          <p className="text-muted-foreground">
            Manage your account and application preferences
          </p>
        </section>

        <Tabs defaultValue="account" className="space-y-4">
          <TabsList>
            <TabsTrigger value="account">Account</TabsTrigger>
            <TabsTrigger value="preferences">Preferences</TabsTrigger>
            <TabsTrigger value="notifications">Notifications</TabsTrigger>
            <TabsTrigger value="advanced">Advanced</TabsTrigger>
            <TabsTrigger value="about">About</TabsTrigger>
          </TabsList>

          <TabsContent value="account" className="space-y-4">
            <AccountSettings />
          </TabsContent>

          <TabsContent value="preferences" className="space-y-4">
            <PreferencesSettings />
          </TabsContent>

          <TabsContent value="notifications" className="space-y-4">
            <NotificationSettings />
          </TabsContent>

          <TabsContent value="advanced" className="space-y-4">
            <AdvancedSettings />
          </TabsContent>

          <TabsContent value="about" className="space-y-4">
            <AboutSettings />
          </TabsContent>
        </Tabs>
      </div>
    </MainLayout>
  );
}

function AccountSettings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Account Settings</CardTitle>
        <CardDescription>
          Manage your personal information and account options
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input id="name" defaultValue="Ashpak Veetar" />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" defaultValue="ashpakv88@gmail.com" />
          </div>
        </div>

        <Separator />

        <div className="space-y-4">
          <h3 className="text-lg font-medium">Danger Zone</h3>

          <div className="grid gap-2">
            <div className="flex flex-col gap-1.5">
              <Label>Export Data</Label>
              <p className="text-sm text-muted-foreground">
                Download all your data including files, folders, and settings
              </p>
            </div>
            <Button variant="outline">Export All Data</Button>
          </div>

          <div className="grid gap-2">
            <div className="flex flex-col gap-1.5">
              <Label>Delete Account</Label>
              <p className="text-sm text-muted-foreground">
                Permanently delete your account and all associated data
              </p>
            </div>
            <Button variant="destructive">Delete Account</Button>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button>Save Changes</Button>
      </CardFooter>
    </Card>
  );
}

function PreferencesSettings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Application Preferences</CardTitle>
        <CardDescription>Customize how Cloudmesh works for you</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="grid gap-2">
            <Label htmlFor="view-mode">Default View Mode</Label>
            <Select defaultValue="list">
              <SelectTrigger id="view-mode">
                <SelectValue placeholder="Select view mode" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="list">List View</SelectItem>
                <SelectItem value="grid">Grid View</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-sm text-muted-foreground">
              Choose how files are displayed by default
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="sort-by">Default Sorting</Label>
            <Select defaultValue="name-asc">
              <SelectTrigger id="sort-by">
                <SelectValue placeholder="Select sorting" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="name-asc">Name (A-Z)</SelectItem>
                <SelectItem value="name-desc">Name (Z-A)</SelectItem>
                <SelectItem value="date-desc">Date (Newest First)</SelectItem>
                <SelectItem value="date-asc">Date (Oldest First)</SelectItem>
                <SelectItem value="size-desc">Size (Largest First)</SelectItem>
                <SelectItem value="size-asc">Size (Smallest First)</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="items-per-page">Items Per Page</Label>
            <Select defaultValue="25">
              <SelectTrigger id="items-per-page">
                <SelectValue placeholder="Select items per page" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="10">10 items</SelectItem>
                <SelectItem value="25">25 items</SelectItem>
                <SelectItem value="50">50 items</SelectItem>
                <SelectItem value="100">100 items</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="landing-page">Default Landing Page</Label>
            <Select defaultValue="dashboard">
              <SelectTrigger id="landing-page">
                <SelectValue placeholder="Select landing page" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="dashboard">Dashboard</SelectItem>
                <SelectItem value="files">File Browser</SelectItem>
                <SelectItem value="transfers">Transfers</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-sm text-muted-foreground">
              Choose which page to show after login
            </p>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button variant="outline">Reset to Defaults</Button>
        <Button>Save Changes</Button>
      </CardFooter>
    </Card>
  );
}

function NotificationSettings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Notification Settings</CardTitle>
        <CardDescription>
          Control when and how you receive notifications
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="transfer-notifications">
                Transfer Notifications
              </Label>
              <p className="text-sm text-muted-foreground">
                Receive notifications when file transfers complete
              </p>
            </div>
            <Switch id="transfer-notifications" defaultChecked />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="quota-alerts">Storage Quota Alerts</Label>
              <p className="text-sm text-muted-foreground">
                Get notified when your storage is almost full
              </p>
            </div>
            <Switch id="quota-alerts" defaultChecked />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="security-notifications">
                Security Notifications
              </Label>
              <p className="text-sm text-muted-foreground">
                Receive alerts about new logins and security events
              </p>
            </div>
            <Switch id="security-notifications" defaultChecked />
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button>Save Changes</Button>
      </CardFooter>
    </Card>
  );
}

function AdvancedSettings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Advanced Settings</CardTitle>
        <CardDescription>
          Configure technical aspects of the application
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="cache-management">Cache Management</Label>
              <p className="text-sm text-muted-foreground">
                Clear cached data to free up space
              </p>
            </div>
            <Button variant="outline" size="sm">
              Clear Cache
            </Button>
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="debug-mode">Debug Mode</Label>
              <p className="text-sm text-muted-foreground">
                Enable detailed logging for troubleshooting
              </p>
            </div>
            <Switch id="debug-mode" />
          </div>

          <Separator />

          <div className="space-y-2">
            <Label>API Usage Statistics</Label>
            <div className="rounded-md border p-4">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>API Calls Today:</span>
                  <span className="font-medium">1,248</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Monthly Limit:</span>
                  <span className="font-medium">50,000</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Reset Date:</span>
                  <span className="font-medium">September 1, 2023</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button variant="outline">Reset to Defaults</Button>
        <Button>Save Changes</Button>
      </CardFooter>
    </Card>
  );
}

function AboutSettings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>About Cloudmesh</CardTitle>
        <CardDescription>Information about the application</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="space-y-1">
            <Label>Version</Label>
            <p className="text-sm">Cloudmesh v1.2.0</p>
          </div>

          <div className="space-y-1">
            <Label>Documentation</Label>
            <div className="flex gap-4">
              <Button variant="link" className="h-auto p-0">
                User Guide
              </Button>
              <Button variant="link" className="h-auto p-0">
                API Documentation
              </Button>
              <Button variant="link" className="h-auto p-0">
                FAQ
              </Button>
            </div>
          </div>

          <div className="space-y-1">
            <Label>Legal</Label>
            <div className="flex gap-4">
              <Button variant="link" className="h-auto p-0">
                Terms of Service
              </Button>
              <Button variant="link" className="h-auto p-0">
                Privacy Policy
              </Button>
            </div>
          </div>

          <div className="space-y-1">
            <Label>Support</Label>
            <p className="text-sm">
              Need help? Contact us at{" "}
              <Button variant="link" className="h-auto p-0">
                support@cloudmesh.com
              </Button>
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
